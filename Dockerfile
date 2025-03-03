# Build stage
FROM golang:1.23-alpine AS builder

# Install required dependencies
RUN apk add --no-cache build-base gcc musl-dev vips-dev

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Install Imaginary
RUN go install github.com/h2non/imaginary@latest

# Build the application
RUN go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install VIPS (required for Imaginary)
RUN apk add --no-cache vips

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary

# Expose necessary ports
EXPOSE 8080 9000

# Run Imaginary as a background process and then start your app
CMD imaginary & ./main
