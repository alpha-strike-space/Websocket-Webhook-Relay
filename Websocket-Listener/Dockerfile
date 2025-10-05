# Minimal Go Image
FROM golang:1.22.4 AS builder
# Working directory
WORKDIR /app
# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app .
# Use a scratch image for a minimal final image
FROM alpine:latest
# Right here in root
WORKDIR /
# Copy the built binary from the builder stage
COPY --from=builder /go-app /go-app
# Run the binary
CMD ["/go-app"]
