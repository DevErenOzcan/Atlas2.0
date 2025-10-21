# Build aşaması
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Bağımlılıkları indirmek için go.mod ve go.sum'ı kopyala
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# go.sum güncellemeleri için tekrar indir
RUN go mod download

# Servisi derle
RUN go build -o /order_service ./main.go

# Çalışma aşaması
FROM alpine:latest

# Gerekli sertifika paketini ekle
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Derlenmiş binary'yi kopyala
COPY --from=builder /order_service .

# Portu açığa çıkar
EXPOSE 50053

# Servisi başlat
CMD ["./order_service"]
