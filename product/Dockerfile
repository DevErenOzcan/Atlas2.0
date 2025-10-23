# product/Dockerfile
# Build aşaması
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Tüm kaynak kodunu kopyala
COPY . .

RUN go mod download
RUN go mod tidy

# Servisi derle
RUN go build -o /product_service ./main.go

# Çalışma aşaması
FROM alpine:latest

# Gerekli sertifika paketini ekle
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Derlenmiş binary'yi kopyala
COPY --from=builder /product_service .

# Portu açığa çıkar
EXPOSE 50052

# Servisi başlat
CMD ["./product_service"]
