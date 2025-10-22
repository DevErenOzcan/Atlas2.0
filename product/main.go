// product/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	// PostgreSQL sürücüsü
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	_ "google.golang.org/protobuf/types/known/timestamppb"

	// Protokol dosyalarından oluşturulan paket
	pb "product/proto"
)

// Server yapısı, gRPC metodlarını implemente eder
type server struct {
	pb.UnimplementedProductServiceServer
	db *sql.DB
}

// ListProducts gRPC metodu
func (s *server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	log.Println("ListProducts isteği alındı.")

	query := `
		SELECT 
			"ProductID", 
			"ProductName", 
			"Description", 
			"Stock", 
			"Price", 
			"Currency",
			"SellerID",
			"Category",
			"DimensDetails"
		FROM "Product" 
		WHERE "IsDeleted" = FALSE`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Veritabanı sorgusu hatasıı: %v", err)
		return nil, fmt.Errorf("ürünler listelenemedi: %w", err)
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var p pb.Product
		var price float64 // Veritabanından NUMERIC'i okumak için

		err := rows.Scan(
			&p.ProductId,
			&p.ProductName,
			&p.Description,
			&p.Stock,
			&price, // Fiyatı float64 olarak oku
			&p.Currency,
			&p.SellerId,
			&p.CategoryId,
			&p.DimensDetails,
		)
		if err != nil {
			log.Printf("Satır okuma hatası: %v", err)
			return nil, fmt.Errorf("ürün verisi okunamadı: %w", err)
		}
		p.Price = price // Okunan fiyatı protobuf mesajına ata

		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Satır döngüsü hatası: %v", err)
		return nil, fmt.Errorf("satır işleme hatası: %w", err)
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func main() {
	// Veritabanı bağlantı dizesini ortam değişkeninden al
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL ortam değişkeni ayarlanmadı")
	}

	// Veritabanı bağlantısı
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}
	defer db.Close()

	// Bağlantıyı kontrol et
	if err = db.Ping(); err != nil {
		log.Fatalf("Veritabanı bağlantısı başarısız: %v", err)
	}
	log.Println("PostgreSQL'e başarıyla bağlanıldı!")

	// gRPC sunucusunu başlat
	port := getEnv("PORT", ":50052")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("TCP dinlenemedi: %v", err)
	}
	log.Printf("Sunucu %s adresinde dinleniyor", port)

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{db: db})

	// Sunucuyu başlat
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hizmet veremedi: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
