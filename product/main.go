// product/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	// PostgreSQL sürücüsü
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	_ "google.golang.org/protobuf/types/known/timestamppb"

	// Prometheus metrics
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "product/proto"
)

type server struct {
	pb.UnimplementedProductServiceServer
	db *sql.DB
}

func (s *server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
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
		log.Printf("Veritabani sorgusu hatasi: %v", err)
		return nil, fmt.Errorf("ürünler listelenemedi: %w", err)
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var p pb.Product
		var price float64

		err := rows.Scan(
			&p.ProductId,
			&p.ProductName,
			&p.Description,
			&p.Stock,
			&price,
			&p.Currency,
			&p.SellerId,
			&p.CategoryId,
			&p.DimensDetails,
		)
		if err != nil {
			log.Printf("Satir okuma hatasi: %v", err)
			return nil, fmt.Errorf("ürün verisi okunamadı: %w", err)
		}
		p.Price = price

		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Satır döngüsü hatası: %v", err)
		return nil, fmt.Errorf("satır işleme hatası: %w", err)
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
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
		WHERE "ProductID" = $1 AND "IsDeleted" = FALSE`

	var p pb.Product
	var price float64

	err := s.db.QueryRowContext(ctx, query, req.GetProductId()).Scan(
		&p.ProductId,
		&p.ProductName,
		&p.Description,
		&p.Stock,
		&price,
		&p.Currency,
		&p.SellerId,
		&p.CategoryId,
		&p.DimensDetails,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ürün bulunamadı")
	} else if err != nil {
		log.Printf("Veritabanı hatası: %v", err)
		return nil, fmt.Errorf("ürün getirilemedi: %w", err)
	}
	p.Price = price

	return &pb.GetProductResponse{Product: &p}, nil
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	query := `
		INSERT INTO "Product" (
			"ProductName", "Description", "Stock", "Price", "Currency", "Alt",
			"SellerID", "Category", "DimensDetails", "IsDeleted"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, FALSE)
		RETURNING "ProductID"`

	var productID int64
	err := s.db.QueryRowContext(ctx, query,
		req.GetProductName(),
		req.GetDescription(),
		req.GetStock(),
		req.GetPrice(),
		req.GetCurrency(),
		req.GetCurrency(),
		req.GetSellerId(),
		req.GetCategoryId(),
		req.GetDimensDetails(),
	).Scan(&productID)

	if err != nil {
		log.Printf("Ürün oluşturma hatası: %v", err)
		return nil, fmt.Errorf("ürün oluşturulamadı: %w", err)
	}

	return &pb.CreateProductResponse{
		ProductId: productID,
		Message:   "Ürün başarıyla oluşturuldu",
	}, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	query := `
		UPDATE "Product"
		SET
			"ProductName" = $1,
			"Description" = $2,
			"Stock" = $3,
			"Price" = $4,
			"Currency" = $5,
			"Alt" = $6,
			"Category" = $7,
			"DimensDetails" = $8
		WHERE "ProductID" = $9 AND "IsDeleted" = FALSE`

	result, err := s.db.ExecContext(ctx, query,
		req.GetProductName(),
		req.GetDescription(),
		req.GetStock(),
		req.GetPrice(),
		req.GetCurrency(),
		req.GetCurrency(),
		req.GetCategoryId(),
		req.GetDimensDetails(),
		req.GetProductId(),
	)

	if err != nil {
		log.Printf("Ürün güncelleme hatası: %v", err)
		return nil, fmt.Errorf("ürün güncellenemedi: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("ürün bulunamadı veya zaten silinmiş")
	}

	return &pb.UpdateProductResponse{
		Message: "Ürün başarıyla güncellendi",
	}, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	query := `
		UPDATE "Product"
		SET "IsDeleted" = TRUE, "DeleteDate" = NOW()
		WHERE "ProductID" = $1 AND "IsDeleted" = FALSE`

	result, err := s.db.ExecContext(ctx, query, req.GetProductId())
	if err != nil {
		log.Printf("Ürün silme hatası: %v", err)
		return nil, fmt.Errorf("ürün silinemedi: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("ürün bulunamadı veya zaten silinmiş")
	}

	return &pb.DeleteProductResponse{
		Message: "Ürün başarıyla silindi",
	}, nil
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL ortam değişkeni ayarlanmadı")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Veritabanı bağlantısı başarısız: %v", err)
	}
	log.Println("PostgreSQL'e başarıyla bağlanıldı!")

	port := getEnv("PORT", ":50052")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("TCP dinlenemedi: %v", err)
	}
	log.Printf("Sunucu %s adresinde dinleniyor", port)

	// gRPC sunucusunu Prometheus için enstrümante et
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterProductServiceServer(s, &server{db: db})
	grpc_prometheus.Register(s)

	// Prometheus /metrics HTTP endpoint'ini ayrı bir portta aç
	metricsPort := getEnv("METRICS_PORT", ":9090")
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("Prometheus metrics %s/metrics adresinde yayında", metricsPort)
		if err := http.ListenAndServe(metricsPort, mux); err != nil {
			log.Fatalf("Metrics sunucusu başlatılamadı: %v", err)
		}
	}()

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
