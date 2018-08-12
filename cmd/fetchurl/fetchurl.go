package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/mholt/archiver"
)

type Template struct {
	To     string
	Name   string
	Url    string
	Sha256 string
	Unpack bool
}

const TEMPLATE string = `
# Add the following snippet to your Dockerfile:
FROM moretea/docker-fetchurl AS {{.Name}}_fetcher
RUN ["fetchurl", \
    "-url", "{{.Url}}", \
    "-sha256", "{{.Sha256}}", \{{if .Unpack}}
    "-unpack"  \ {{end}}
    "-to", "{{.To}}"]

# And use in another layer like:
FROM ...
...
COPY --from={{.Name}}_fetcher {{.To}} {{.To}}
`

func main() {
	var rawURL string
	var expectedSha256 string
	var to string
	var doTemplate bool
	var doUnpack bool

	flag.StringVar(&rawURL, "url", "", "URL to download")
	flag.StringVar(&expectedSha256, "sha256", "", "The sha256 of the to-download file")
	flag.StringVar(&to, "to", "", "Where to store the end result")
	flag.BoolVar(&doTemplate, "template", false, "Download, compute sha256 and print out usage template")
	flag.BoolVar(&doUnpack, "unpack", false, "Unpack the archive")
	flag.Parse()

	if rawURL == "" {
		fatal("No URL provided; add -url '$URL' to your invocation of fetchurl")
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		fatal("URL is invalid; %v", err)
	}

	if !doTemplate && to == "" {
		fatal("No target location is specified; add -to or -template to your invocation of fetchurl")
	}

	if !doTemplate && expectedSha256 == "" {
		fatal("I can't download unverified files. Please add -sha256 to your invocation of fetchurl.\nNote: you can find the sha256 by running fetchurl -url $MY_URL -template")
	}

	fmt.Printf("Downloading '%s'...", url)
	file, err := download(url)
	if err != nil {
		fatal("Could not download file, because: %v\n", err)
	}
	fmt.Printf(" Done!\n")

	defer os.Remove(file.Name())

	computedSha256, err := computeSha256(file)
	if err != nil {
		fatal("Could not compute SHA256 of the downloaded file, because: %v\n", err)
	}

	if doTemplate {
		printTemplate(url, computedSha256)
	} else {
		if expectedSha256 != computedSha256 {
			fatal("Hmm. I expected a sha256 of '%s', but instead computed '%s'", expectedSha256, computedSha256)
		}

		if doUnpack {
			format := archiver.MatchingFormat(file.Name())
			if format == nil {
				fatal("Could not unpack the downloaded file; not a supported archive")
			}

			err = format.Open(file.Name(), to)
			if err != nil {
				fatal("Could not unpack the downloaded file; %v", err)
			}

		} else {
			err = copyFile(file, to)
			if err != nil {
				fatal("Could not rename the downloaded file to %s, because: %v", to, err)
			}
			fmt.Printf("Downloaded to %s\n", to)
		}
	}
}

func printTemplate(url *url.URL, sha256 string) {
	to := "/"

	if url.Path == "" {
		host := strings.Split(url.Host, ":")[0]
		to += host
	} else {
		parts := strings.Split(url.Path, "/")
		lastPart := parts[len(parts)-1]
		to += lastPart
	}

	// Create a name for the fetcher layer;
	// 1. Remove al non-letter names from the start of the name
	// 2. Replace all (consecutive) non ascii letters from the name
	re := regexp.MustCompile("^[^a-z]+")
	name := re.ReplaceAllString(strings.ToLower(to), "")
	re = regexp.MustCompile("[^a-z]+")
	name = re.ReplaceAllString(name, "_")

	// Now perform an heuristic check to determine whether or not we should suggest the end user to unpack the file.
	var archiveSuffixes = []string{
		".zip",
		".tar",
		".tar.gz", ".tgz",
		".tar.bz2", ".tbz2",
		".tar.xz", ".txz",
		".tar.lz4", ".tlz4",
		".tar.sz", ".tsz",
		".rar",
	}

	unpack := false

	for _, archiveSuffix := range archiveSuffixes {
		if strings.HasSuffix(to, archiveSuffix) {
			unpack = true
			break
		}
	}

	dockerfileTemplate := Template{Name: name, Url: url.String(), To: to, Sha256: sha256, Unpack: unpack}

	tmpl, err := template.New("").Parse(TEMPLATE)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, dockerfileTemplate)
	if err != nil {
		panic(err)
	}
}

func download(url *url.URL) (file *os.File, err error) {
	file, err = ioutil.TempFile("", "fetchurl")
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.Get(url.String())
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Download unsuccesful; status code %v", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return
}

func computeSha256(file *os.File) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, file)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func copyFile(from *os.File, toPath string) error {
	from.Seek(0, 0)
	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(to, from)
	return err
}

func fatal(format string, data ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", data...)
	os.Exit(1)
}
