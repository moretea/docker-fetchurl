# Fetchurl Docker Image

Ever annoyed by the fact that downloading files, and verifying that they are correct, is pretty annoying in Docker?

This image helps you download URL's, verify their content, and cache them so that you don't re-download them all the time!

The idea & name are inspired by [nixpkg](https://nixos.org/nixpkgs/)'s fetchurl function.

## Usage

You can use this image as shown below:

```Dockerfile
ARG X_VERSION 1.0
ARG X_SHA256 8a51c03f1ff77c2b8e76da512070c23c5e69813d5c61732b3025199e5f0c14d5

FROM moretea/fetchurl AS download_x
RUN fetchurl \
  --url http://example.com/download/x-${X_VERSION}.tar.gz \
  --sha256 ${X_SHA256} \
  --to /x.tar.gz

FROM alpine
COPY --from=download_x /x /opt/x
ENTRYPOINT ["/opt/x/bin/x"]
```

If you want to know what SHA256 an image has, run:

```
$ docker run --rm moretea/fetchurl http://example.com/download/x-1.0.tar.gz
```

## Example
Also checkout [Dockerfile.example](./Dockerfile.example), where we download example.com and serve our own copy with Nginx.


```
$ docker build -t example.com -f Dockerfile.example 
$ docker run --rm -ti -p 8080:80 example.com
```
