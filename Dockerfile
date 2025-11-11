# Stage 1: Build
FROM golang:1.25.1-alpine AS builder

# Install git untuk dependency Go module
RUN apk add --no-cache git

# Set working directory di dalam container
WORKDIR /app

# Copy module files dan download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh project ke container
COPY . .

# Build aplikasi Go
RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

# Working directory untuk runtime
WORKDIR /root/

# Copy hasil build dari stage pertama
COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads

# Expose port (Railway akan otomatis override)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
