// user/main.go
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
	pb "user/proto"
)

// userRoleMap, PostgreSQL ENUM'u ile Protobuf ENUM'u arasındaki eşleşmeyi sağlar
var userRoleMap = map[string]pb.UserRole{
	"KULLANICI": pb.UserRole_KULLANICI,
	"SATICI":    pb.UserRole_SATICI,
	"YONETICI":  pb.UserRole_YONETICI,
}

// Server yapısı, gRPC metodlarını implemente eder
type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

// UserRoleFromDB, veritabanından gelen string rolü Protobuf enum'a dönüştürür
func UserRoleFromDB(role string) pb.UserRole {
	if r, ok := userRoleMap[role]; ok {
		return r
	}
	return pb.UserRole_KULLANICI // Varsayılan rol
}

// ListUsers gRPC metodu
func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Println("ListUsers isteği alındı")

	// Sadece User tablosundan gerekli sütunları çek
	query := `
		SELECT 
			"UserID", 
			"Email", 
			"Ad", 
			"Soyad", 
			"Role", 
			"AddressID"
		FROM "User" 
		WHERE "IsDeleted" = FALSE`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Veritabanı sorgusu hatası: %v", err)
		return nil, fmt.Errorf("kullanıcılar listelenemedi: %w", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var (
			userID    int64
			email     string
			ad        sql.NullString
			soyad     sql.NullString
			role      string // Veritabanından string olarak al
			addressID sql.NullInt64
		)

		err := rows.Scan(&userID, &email, &ad, &soyad, &role, &addressID)
		if err != nil {
			log.Printf("Satır okuma hatası: %v", err)
			return nil, fmt.Errorf("kullanıcı verisi okunamadı: %w", err)
		}

		user := &pb.User{
			UserId: userID,
			Email:  email,
			Ad:     ad.String,
			Soyad:  soyad.String,
			Role:   UserRoleFromDB(role),
			// AddressID NULL olabilir, kontrol et
			AddressId: addressID.Int64,
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Satır döngüsü hata: %v", err)
		return nil, fmt.Errorf("satır işleme hatası: %w", err)
	}

	return &pb.ListUsersResponse{Users: users}, nil
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
	port := getEnv("PORT", ":50051")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("TCP dinlenemedi: %v", err)
	}
	log.Printf("Sunucu %s adresinde dinleniyor", port)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	// Sunucuyu başlat
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hizmet vermedi: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
