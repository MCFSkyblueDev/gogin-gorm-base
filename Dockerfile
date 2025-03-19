# STEP 1: Build the Go application
FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# STEP 2: Create a minimal final image
FROM alpine:3.20

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

# Copy static files from the uploads folder
COPY uploads ./uploads

EXPOSE 8080

CMD ["./main"]