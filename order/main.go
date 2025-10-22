// order/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "order/proto"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log.Println("ListOrders isteği alındı")

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

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{db: db})

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
