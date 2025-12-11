# Stage 1: Builder (Image besar untuk compile)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files dulu agar cache layer optimal
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary aplikasi
# -o main: nama output binary
# cmd/main.go: lokasi entry point
RUN go build -o main cmd/main.go

# Stage 2: Runner (Image kecil untuk produksi)
FROM alpine:latest

WORKDIR /root/

# Install sertifikat SSL (penting jika nembak API HTTPS luar)
RUN apk --no-cache add ca-certificates

# Copy binary dari Stage 1
COPY --from=builder /app/main .
# Copy folder migrasi (penting agar migrate jalan)
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Expose port yang dipakai aplikasi
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["./main"]