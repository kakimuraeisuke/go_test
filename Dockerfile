# Build stage
FROM golang:1.22-alpine AS builder

# Install protoc and required tools
RUN apk add --no-cache \
    protobuf \
    protobuf-dev \
    git \
    make

# Install protoc-gen-go and protoc-gen-go-grpc (compatible with Go 1.22)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Install grpc-health-probe (compatible with Go 1.22)
RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.24

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf files
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/go_test/v1/go_test.proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy grpc-health-probe
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/

# Set working directory
WORKDIR /root/

# Copy the binary
COPY --from=builder /app/main .

# Expose port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD grpc-health-probe -addr=:50051 || exit 1

# Run the application
CMD ["./main"]