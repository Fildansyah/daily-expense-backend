# Stage 1: Build binary dengan Go 1.24
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy module files dan download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code dan build
# Catatan: main package ada di cmd/main.go (sesuai makefile), bukan di root
COPY . .
RUN go build -o main ./cmd

# Stage 2: Minimal image untuk running
FROM alpine:latest
WORKDIR /app

# ca-certificates dibutuhkan jika app memanggil API HTTPS eksternal
RUN apk add --no-cache ca-certificates

# Copy binary dari stage builder
COPY --from=builder /app/main .

EXPOSE 8000
CMD ["./main"]
