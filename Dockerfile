# Use the official Go base image
FROM docker.io/library/golang:alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files into the container
# only do this once a go.sum exists, of course
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the remaining Go source code into the container
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a minimal alpine image with Caddy
FROM docker.io/library/caddy:latest

# Copy the Go binary from the builder stage
COPY --from=builder /app/app /app

# Copy assets files
COPY assets /assets

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Expose port 8080 for the application
EXPOSE 8080

# Start the Go application and Caddy
CMD ["/bin/sh", "-c", "/app --server & /usr/bin/caddy run --config /etc/caddy/Caddyfile"]
