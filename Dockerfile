# FROM golang:1.21-alpine
FROM golang:1.25-alpine

# Set workdir
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o email-service main.go handler.go worker.go queue.go model.go

# Expose port 8080
EXPOSE 8080

# Default environment variables (can override)
ENV QUEUE_SIZE=10
ENV NUM_WORKERS=3

# Run the binary
CMD ["./email-service"]
