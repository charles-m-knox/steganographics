FROM docker.io/library/golang:alpine AS builder

RUN apk add upx

RUN mkdir -p /site/pkg/steganographics
WORKDIR /site
COPY go.mod /site
COPY pkg/steganographics/go.mod /site/pkg/steganographics/go.mod

RUN go mod download

COPY . /site

WORKDIR /site
RUN CGO_ENABLED=0 go build -v -o bin-uncompressed -ldflags="-w -s -buildid=" -trimpath
RUN upx --best -o ./steganographics bin-uncompressed

FROM docker.io/library/alpine:latest
COPY --from=builder /site/steganographics /steganographics

LABEL org.opencontainers.image.source https://github.com/charles-m-knox/steganographics

CMD ["/steganographics"]
