# Fetchurl Docker Image

<a href="https://hub.docker.com/r/moretea/docker-fetchurl/"><img alt="Image size badge" src="https://img.shields.io/microbadger/image-size/moretea/docker-fetchurl.svg"/> <img alt="Docker Build Status" src="https://img.shields.io/docker/build/moretea/docker-fetchurl.svg"/></a>

_Ever annoyed by all the hoops you have to jump through if you want to download and verify that you downloaded the correct file in a Dockerfile?_

This image helps you download URL's, verify their content, unpack them if they are an archive, and cache them so that you don't have to re-download them all the time!

The idea & name are inspired by [nixpkg](https://nixos.org/nixpkgs/)'s fetchurl function.

## Quick start
1. Run `docker run --rm moretea/docker-fetchurl -template -url $MY_URL_TO_DOWNLOAD`
2. Copy & paste the output in your Dockerfile.

```
$ docker run --rm moretea/docker-fetchurl http://maarten-hoogendoorn.nl/blog
Downloading 'http://maarten-hoogendoorn.nl/blog'... Done!

# Add the following snippet to your Dockerfile:
FROM moretea/docker-fetchurl AS blog_fetcher
RUN ["fetchurl", \
    "-url", "http://maarten-hoogendoorn.nl/blog", \
    "-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", \
    "-to", "/blog"]

# And use in another layer like:
FROM ...
...
COPY --from=blog_fetcher /blog /blog
```

## Usage

```
$ docker run moretea/docker-fetchurl --help
Usage of /bin/fetchurl:
  -sha256 string
    	The sha256 of the to-download file
  -template
    	Download, compute sha256 and print out usage template
  -to string
    	Where to store the end result
  -unpack
    	Unpack the archive
  -url string
    	URL to download
```

## Serving example.com
Check [Dockerfile.example](./Dockerfile.example) for an example where we download the content of example.com and serve it with our own nginx instance.

Build your own example.com enginx container with:

```
$ docker build -t example.com https://raw.githubusercontent.com/moretea/docker-fetchurl/master/Dockerfile.example
$ docker run -d -p 8080:80 example.com
$ curl localhost:8080
  <html>
  ....
  </html>
```
