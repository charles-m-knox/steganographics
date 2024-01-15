# steganographics

A tool to enable easy access to steganography.

## Usage

This program has two modes of operation - either as a command line tool, or as a standalone HTTP server. If you run it using `make run` below, a Docker container will run both Caddy and the application itself simultaneously, with Caddy doing gzip compression and static file hosting.

```bash
make prep
make build
make run

./steganographics --help

# the below two are required at the same time
# caddy serves static files locally
make run-caddy
make run-server
```

When running as an HTTP server, this program provides two REST API endpoints:

1. `/api/hide`: Accepts a POST request with a JSON payload containing a `msg` field with the text to hide, and a `png` field with the Base64-encoded PNG image data. The endpoint hides the text in the image using the LSB technique and returns the modified PNG image.

2. `/api/extract`: Accepts a POST request with a JSON payload containing a `png` field with the Base64-encoded PNG image data. The endpoint extracts the hidden text from the image using the LSB technique and returns a JSON response with the extracted text in a `msg` field.

To test the endpoints, you can use an HTTP client like `curl` or a tool like Postman to send POST requests with the appropriate JSON payload to the `/api/hide` and `/api/extract` endpoints.

## Appendix

Useful command:

```bash
convert -quality 100 -define png:compression-level=9 input.jpg input.png
```
