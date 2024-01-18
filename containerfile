FROM docker.io/library/golang:alpine AS builder

RUN apk add upx

WORKDIR /site
COPY go.mod /site

RUN go mod download

COPY . /site

WORKDIR /site
RUN CGO_ENABLED=0 go build -v -o bin-uncompressed -ldflags="-w -s -buildid=" -trimpath
RUN upx --best -o ./steganographics bin-uncompressed

FROM docker.io/library/alpine:latest
COPY --from=builder /site/steganographics /steganographics

CMD ["/steganographics"]
