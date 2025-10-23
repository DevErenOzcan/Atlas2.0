// order/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	// Prometheus metrics
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "order/proto"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	query := `
		SELECT
			o."OrderID", o."UserID", o."Total", o."Final", o."IsShipped", o."CreateDate",
			od."DetailID", od."ProductID", od."Item", od."Final" as DetailFinal, od."CreateDate" as DetailCreateDate
		FROM "Order" o
		LEFT JOIN "OrderDetail" od ON o."OrderID" = od."OrderID"
		WHERE o."IsDeleted" = FALSE
		ORDER BY o."OrderID" ASC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Veritabanı sorgusu hatası: %v", err)
		return nil, fmt.Errorf("siparişler listelenemedi: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[int64]*pb.Order)
	var orders []*pb.Order

	for rows.Next() {
		var (
			orderID          int64
			userID           int64
			total            float64
			final            float64
			isShipped        bool
			createDate       sql.NullTime
			detailID         sql.NullInt64
			productID        sql.NullInt64
			item             sql.NullInt32
			detailFinal      sql.NullFloat64
			detailCreateDate sql.NullTime
		)

		err := rows.Scan(
			&orderID, &userID, &total, &final, &isShipped, &createDate,
			&detailID, &productID, &item, &detailFinal, &detailCreateDate,
		)
		if err != nil {
			log.Printf("Satır okuma hatası: %v", err)
			return nil, fmt.Errorf("sipariş verisi okunamadı: %w", err)
		}

		if _, ok := ordersMap[orderID]; !ok {
			order := &pb.Order{
				OrderId:      orderID,
				UserId:       userID,
				Total:        total,
				Final:        final,
				IsShipped:    isShipped,
				CreateDate:   timestamppb.New(createDate.Time),
				OrderDetails: []*pb.OrderDetail{},
			}
			ordersMap[orderID] = order
			orders = append(orders, order)
		}

		if detailID.Valid {
			orderDetail := &pb.OrderDetail{
				DetailId:   detailID.Int64,
				OrderId:    orderID,
				ProductId:  productID.Int64,
				Item:       item.Int32,
				Final:      detailFinal.Float64,
				CreateDate: timestamppb.New(detailCreateDate.Time),
			}
			ordersMap[orderID].OrderDetails = append(ordersMap[orderID].OrderDetails, orderDetail)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Satır döngüsü hatası: %v", err)
		return nil, fmt.Errorf("satır işleme hatası: %w", err)
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Transaction başlatılamadı: %v", err)
		return nil, fmt.Errorf("sipariş oluşturulamadı: %w", err)
	}
	defer tx.Rollback()

	var total, final float64
	for _, item := range req.GetItems() {
		itemTotal := float64(item.GetQuantity()) * item.GetPrice()
		total += itemTotal
		final += itemTotal
	}

	var orderID int64
	orderQuery := `
		INSERT INTO "Order" ("UserID", "Total", "Final", "IsShipped", "IsDeleted")
		VALUES ($1, $2, $3, FALSE, FALSE)
		RETURNING "OrderID"`

	err = tx.QueryRowContext(ctx, orderQuery, req.GetUserId(), total, final).Scan(&orderID)
	if err != nil {
		log.Printf("Order oluşturma hatası: %v", err)
		return nil, fmt.Errorf("sipariş kaydedilemedi: %w", err)
	}

	detailQuery := `
		INSERT INTO "OrderDetail" ("OrderID", "ProductID", "Item", "Final", "IsDeleted")
		VALUES ($1, $2, $3, $4, FALSE)`

	for _, item := range req.GetItems() {
		itemFinal := float64(item.GetQuantity()) * item.GetPrice()
		_, err = tx.ExecContext(ctx, detailQuery, orderID, item.GetProductId(), item.GetQuantity(), itemFinal)
		if err != nil {
			log.Printf("OrderDetail oluşturma hatası: %v", err)
			return nil, fmt.Errorf("sipariş detayı kaydedilemedi: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit hatası: %v", err)
		return nil, fmt.Errorf("sipariş tamamlanamadı: %w", err)
	}

	return &pb.CreateOrderResponse{
		OrderId: orderID,
		Message: "Sipariş başarıyla oluşturuldu",
	}, nil
}

func (s *server) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	query := `
		UPDATE "Order"
		SET "IsShipped" = $1
		WHERE "OrderID" = $2 AND "IsDeleted" = FALSE`

	result, err := s.db.ExecContext(ctx, query, req.GetIsShipped(), req.GetOrderId())
	if err != nil {
		log.Printf("Sipariş durumu güncelleme hatası: %v", err)
		return nil, fmt.Errorf("sipariş durumu güncellenemedi: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("sipariş bulunamadı veya zaten silinmiş")
	}

	statusText := "onaylandı"
	if !req.GetIsShipped() {
		statusText = "iptal edildi"
	}

	return &pb.UpdateOrderStatusResponse{
		Message: fmt.Sprintf("Sipariş başarıyla %s", statusText),
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

	port := getEnv("PORT", ":50053")
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
	pb.RegisterOrderServiceServer(s, &server{db: db})
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
