# Build aşaması
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Tüm kaynak kodunu kopyala (proto dahil)
COPY . .

RUN go mod download
RUN go mod tidy

# Derleme
RUN go build -o /gateway_service ./main.go

# Çalışma aşaması
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /gateway_service .
EXPOSE 8080
CMD ["./gateway_service"]
