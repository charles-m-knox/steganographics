# steganographics

A tool to enable easy access to steganography by embedding secret text in PNG images.

It uses a simple Least Significant Bit (LSB) algorithm to distribute the desired secret bytes across the first bytes in the red channel of an image.

## Usage

This program has three modes of operation:

- as a Go library
- command line tool
- a standalone HTTP server.

### Go library

This Go library has no dependencies aside from the standard library, and will work with `CGO_ENABLED=0`.

But first, you must install it like this:

```bash
GOSUMDB=off GOPROXY=https://gitea.cmcode.dev/api/packages/cmcode/go go get -v gitea.cmcode.dev/cmcode/steganographics/secrets@latest
```

> *If you do not do this, you will not be able to use the library - Go's built-in package resolver will attempt to use `main.go` and will fail.*

Usage:

```go
package main

import (
    s "gitea.cmcode.dev/cmcode/steganographics/secrets"
)

func getTextFromStdin() {
    img, _, err := image.Decode(os.Stdin)
    if err != nil {
        log.Fatalf("failed to read img from stdin: %v", err.Error())
    }

    os.Stdout.Write(s.ExtractTextFromImage(img))
}

func writeTextToImageFromStdin(hiddenText string) {
    img, _, err := image.Decode(os.Stdin)
    if err != nil {
        log.Fatalf("failed to read img from stdin: %v", err.Error())
    }

    output, err := s.HideTextInImage(img, []byte(hiddenText))
    if err != nil {
        log.Fatalf("failed to read img from stdin: %v", err.Error())
    }

    err = png.Encode(os.Stdout, output)
    if err != nil {
        log.Fatalf("failed to write img to stdout: %v", err.Error())
    }
}
```

### Command line tool

If you want to use Go's `go install` to install steganographics, you can do it by appending `-cli` to the version:

```bash
GOSUMDB=off GOPROXY=https://gitea.cmcode.dev/api/packages/cmcode/go go install gitea.cmcode.dev/cmcode/steganographics@v0.1.1-cli
```

> *Note: If you forget to add the suffix `-cli`, it will instead install the library-only version.*

Once installed, `steganographics` supports piping from `stdin` and to `stdout`:

```bash
# writes secret text to an image, and then immediately echoes the secret
# text from it
cat path/to/image.png | ./steganographics -secret "secret message1" | ./steganographics
```

It can also directly read and write to/from files:

```bash
# write a message into an image:
./steganographics -input path/to/image.png -output path/to/output.png -secret "secret message2"

# read a message from an image:
./steganographics -input path/to/image.png
```

### HTTP server API

```bash
./steganographics -addr "0.0.0.0:29104"

# with tls:
make gen-tls-certs # optional - generates certs and requires interactive input
./steganographics -addr "0.0.0.0:29104" -cert cert.pem -key key.pem
```

When running as an HTTP server, this program provides two REST API endpoints:

1. `/api/hide`: Accepts a POST request with a JSON payload containing a `msg` field with the text to hide, and a `png` field with the Base64-encoded PNG image data. The endpoint hides the text in the image using the LSB technique and returns the modified PNG image.

2. `/api/extract`: Accepts a POST request with a JSON payload containing a `png` field with the Base64-encoded PNG image data. The endpoint extracts the hidden text from the image using the LSB technique and returns a JSON response with the extracted text in a `msg` field.

To test the endpoints, you can use an HTTP client like `curl` or a tool like Postman to send POST requests with the appropriate JSON payload to the `/api/hide` and `/api/extract` endpoints.

## Disclaimers and warnings

- This may not be the most secure/private/hardened method of encoding text in an image.
- This is a side project and may have bugs or other issues.
- This project is licensed AGPL3 so you cannot use it unless you abide by the license terms.
- Unit tests are not currently implemented, but the library does work as expected in intended, normal use cases.
- This library may alter your source image in more ways than simply the LSB algorithm.
- Compression or other image alterations run the risk of wiping out the embedded image.

## Appendix

This is a useful command using ImageMagick's `covert` command line tool to convert a jpg to png:

```bash
convert -quality 100 -define png:compression-level=9 input.jpg input.png
```
