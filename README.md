# Fetchurl Docker Image

Ever annoyed by the fact that downloading files, and verifying that they are correct, is pretty annoying in Docker?

This image helps you download URL's, verify their content, and cache them so that you don't re-download them all the time!

The idea & name are inspired by [nixpkg](https://nixos.org/nixpkgs/)'s fetchurl function.

## Usage
1. Run `docker run --rm moretea/fetchurl $MY_URL_TO_DOWNLOAD`
2. Copy & Paste the output in your Dockerfile.

```
$ docker run --rm moretea/fetchurl http://maarten-hoogendoorn.nl/blog
Downloading 'http://maarten-hoogendoorn.nl/blog'... Done!

# Add the following snippet to your Dockerfile:
FROM moretea/fetchurl AS blog_fetcher
RUN ["fetchurl", \
    "-url", "http://maarten-hoogendoorn.nl/blog", \
    "-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", \
    "-to", "/blog"]

# And use in another layer like:
FROM ...
...
COPY --from=blog /blog /blog
```

## Serving example.com
Check [Dockerfile.example](./Dockerfile.example) for an example where we download the content of example.com and serve it with our own nginx instance.
