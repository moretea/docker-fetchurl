package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var url string
	var sha256 string
	var to string
	var doTest bool

	flag.StringVar(&url, "url", "", "URL to download")
	flag.StringVar(&sha256, "sha256", "", "The sha256 of the to-download file")
	flag.StringVar(&to, "to", "", "Where to store the end result")
	flag.BoolVar(&doTest, "test", false, "Ony download, and print the sha256")
	flag.Parse()

	if url == "" {
		fatal("No URL provided; add --url 'url' to your invocation of fetchurl")
	}

	if !doTest && to == "" {
		fatal("No target location is specified; add --to or --test to your invocation of fetchurl")
	}

	if !doTest && sha256 == "" {
		fatal("I can't download unverified files. Please add --sha256 to your invocation of fetchurl.\nNote: you can find the sha256 by running fetchurl --url $MY_URL --test")
	}

	fmt.Printf("Downloading '%s'...", url)
	fmt.Printf(" Done!\n")
	file, err := download(url)

	if err != nil {
		fatal("Could not download file, because: %v\n", err)
	}

	defer os.Remove(file.Name())

	sha, err := computeSha256(file)

	if err != nil {
		fatal("Could not compute SHA256 of the downloaded file, because: %v\n", err)
	}

	if doTest {
		fmt.Printf(`fetchurl --url "%s" --sha256 %s`, url, sha)
	} else {
		err = copyFile(file, to)
		if err != nil {
			fatal("Could not rename the downloaded file to %s, because: %v", to, err)
		}
		fmt.Printf("Downloaded to %s\n", to)
	}
}

func download(url string) (file *os.File, err error) {
	file, err = ioutil.TempFile("", "fetchurl")
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.Get(url)
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
