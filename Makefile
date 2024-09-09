.PHONY=build

BUILDDIR=build
VER=0.1.5
BIN=$(BUILDDIR)/steganographics-v$(VER)
REPO=github.com/charles-m-knox/steganographics

build-dev:
	CGO_ENABLED=0 go build -ldflags="-X main.version=v$(VER)-dev" -v

mkbuilddir:
	mkdir -p $(BUILDDIR)

build-prod: mkbuilddir
	CGO_ENABLED=0 go build -v -o $(BIN) -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath

test:
	go test -test.v -coverprofile=testcov.out ./... && \
	go tool cover -html=testcov.out

run:
	./$(BIN)

lint:
	golangci-lint run ./...

compress-prod: mkbuilddir
	rm -f $(BIN)-compressed
	upx --best -o ./$(BIN)-compressed $(BIN)

build-mac-arm64: mkbuilddir
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o $(BIN)-darwin-arm64 -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-arm64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-arm64

build-mac-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o $(BIN)-darwin-amd64 -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-amd64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-amd64

build-win-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o $(BIN)-win-amd64-uncompressed -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath
	rm -f $(BIN)-win-amd64
	upx --best -o ./$(BIN)-win-amd64 $(BIN)-win-amd64-uncompressed

build-linux-arm64: mkbuilddir
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o $(BIN)-linux-arm64-uncompressed -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-arm64
	upx --best -o ./$(BIN)-linux-arm64 $(BIN)-linux-arm64-uncompressed

build-linux-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o $(BIN)-linux-amd64-uncompressed -ldflags="-X main.version=v$(VER) -w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-amd64
	upx --best -o ./$(BIN)-linux-amd64 $(BIN)-linux-amd64-uncompressed

build-all: mkbuilddir build-linux-amd64 build-linux-arm64 build-win-amd64 build-mac-amd64 build-mac-arm64

delete-builds:
	rm $(BUILDDIR)/*

gen-tls-certs:
	openssl genrsa -out key.pem 2048 && \
	openssl ecparam -genkey -name secp384r1 -out key.pem && \
	openssl req -new -x509 -sha256 -key key.pem -out cert.pem -days 3650

podman-build:
	podman build -t $(REPO):latest -f containerfile .
	podman tag $(REPO):latest $(REPO):v$(VER)

# requires you to run 'podman login ghcr.io'
podman-push:
	podman push $(REPO):latest
	podman push $(REPO):v$(VER)

podman-run:
	podman rm -f steganographics || true
	podman run -d \
		-p "127.0.0.1:29104:29104" \
		--restart=unless-stopped \
		--name=steganographics \
		-it steganographics:latest \
			/steganographics -addr 0.0.0.0:29104
	podman logs -f steganographics
