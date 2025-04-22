# Use official Golang image as a builder stage
FROM golang:1.24.1 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/web-analyzer .

RUN chmod +x /app/web-analyzer

# Final stage: Use a smaller base image
FROM alpine:latest

# Set working directory
WORKDIR /root/

COPY --from=builder /app /app

# Copy the binary from the builder stage
COPY --from=builder /app/web-analyzer .

RUN apk add --no-cache ca-certificates

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./web-analyzer"]
