FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o app-core ./cmd/app

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app-core .

EXPOSE 8080

CMD ["./app-core"]
