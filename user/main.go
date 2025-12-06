// user/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	// PostgreSQL sürücüsü
	_ "github.com/lib/pq"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "google.golang.org/protobuf/types/known/timestamppb"

	// Prometheus metrics
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Protokol dosyalarından oluşturulan paket
	pb "user/proto"
)

var userRoleMap = map[string]pb.UserRole{
	"KULLANICI": pb.UserRole_KULLANICI,
	"SATICI":    pb.UserRole_SATICI,
	"YONETICI":  pb.UserRole_YONETICI,
}

type server struct {
	pb.UnimplementedUserServiceServer
	db        *sql.DB
	jwtSecret string
}

func UserRoleFromDB(role string) pb.UserRole {
	if r, ok := userRoleMap[role]; ok {
		return r
	}
	return pb.UserRole_KULLANICI
}

func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
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
			role      string
			addressID sql.NullInt64
		)

		err := rows.Scan(&userID, &email, &ad, &soyad, &role, &addressID)
		if err != nil {
			log.Printf("Satir okuma hatasi: %v", err)
			return nil, fmt.Errorf("kullanıcı verisi okunamadı: %w", err)
		}

		user := &pb.User{
			UserId:    userID,
			Email:     email,
			Ad:        ad.String,
			Soyad:     soyad.String,
			Role:      UserRoleFromDB(role),
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

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email ve şifre gereklidir")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Şifre hash hatası: %v", err)
		return nil, status.Error(codes.Internal, "Şifre işlenemedi")
	}

	query := `
		INSERT INTO "User" ("Email", "Password", "Ad", "Soyad", "Role", "IsDeleted")
		VALUES ($1, $2, $3, $4, 'KULLANICI', FALSE)
		RETURNING "UserID"`

	var userID int64
	err = s.db.QueryRowContext(ctx, query, req.GetEmail(), string(hashedPassword), req.GetAd(), req.GetSoyad()).Scan(&userID)
	if err != nil {
		log.Printf("Kullanıcı kayıt hatası: %v", err)
		return nil, status.Error(codes.AlreadyExists, "Bu email zaten kullanılıyor")
	}

	return &pb.RegisterResponse{
		UserId:  userID,
		Email:   req.GetEmail(),
		Message: "Kayıt başarılı",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email ve şifre gereklidir")
	}

	query := `
		SELECT "UserID", "Password", "Role"
		FROM "User"
		WHERE "Email" = $1 AND "IsDeleted" = FALSE`

	var (
		userID       int64
		hashedPasswd string
		role         string
	)

	err := s.db.QueryRowContext(ctx, query, req.GetEmail()).Scan(&userID, &hashedPasswd, &role)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "Kullanıcı bulunamadı")
	} else if err != nil {
		log.Printf("Veritabanı hatası: %v", err)
		return nil, status.Error(codes.Internal, "Giriş yapılamadı")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Hatalı şifre")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   req.GetEmail(),
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Printf("Token oluşturma hatası: %v", err)
		return nil, status.Error(codes.Internal, "Token oluşturulamadı")
	}

	return &pb.LoginResponse{
		Token:  tokenString,
		UserId: userID,
		Email:  req.GetEmail(),
		Role:   UserRoleFromDB(role),
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	query := `
		SELECT "UserID", "Email", "Ad", "Soyad", "Role", "AddressID"
		FROM "User"
		WHERE "UserID" = $1 AND "IsDeleted" = FALSE`

	var (
		userID    int64
		email     string
		ad        sql.NullString
		soyad     sql.NullString
		role      string
		addressID sql.NullInt64
	)

	err := s.db.QueryRowContext(ctx, query, req.GetUserId()).Scan(&userID, &email, &ad, &soyad, &role, &addressID)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "Kullanıcı bulunamadı")
	} else if err != nil {
		log.Printf("Veritabanı hatası: %v", err)
		return nil, status.Error(codes.Internal, "Kullanıcı getirilemedi")
	}

	user := &pb.User{
		UserId:    userID,
		Email:     email,
		Ad:        ad.String,
		Soyad:     soyad.String,
		Role:      UserRoleFromDB(role),
		AddressId: addressID.Int64,
	}

	return &pb.GetUserResponse{User: user}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	roleStr := "KULLANICI"
	switch req.GetRole() {
	case pb.UserRole_SATICI:
		roleStr = "SATICI"
	case pb.UserRole_YONETICI:
		roleStr = "YONETICI"
	}

	query := `
		UPDATE "User"
		SET "Email" = $1, "Ad" = $2, "Soyad" = $3, "Role" = $4
		WHERE "UserID" = $5 AND "IsDeleted" = FALSE`

	result, err := s.db.ExecContext(ctx, query,
		req.GetEmail(),
		req.GetAd(),
		req.GetSoyad(),
		roleStr,
		req.GetUserId(),
	)

	if err != nil {
		log.Printf("Kullanıcı güncelleme hatası: %v", err)
		return nil, status.Error(codes.Internal, "Kullanıcı güncellenemedi")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "Kullanıcı bulunamadı veya zaten silinmiş")
	}

	return &pb.UpdateUserResponse{
		Message: "Kullanıcı başarıyla güncellendi",
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `
		UPDATE "User"
		SET "IsDeleted" = TRUE, "DeleteDate" = NOW()
		WHERE "UserID" = $1 AND "IsDeleted" = FALSE`

	result, err := s.db.ExecContext(ctx, query, req.GetUserId())
	if err != nil {
		log.Printf("Kullanıcı silme hatası: %v", err)
		return nil, status.Error(codes.Internal, "Kullanıcı silinemedi")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "Kullanıcı bulunamadı veya zaten silinmiş")
	}

	return &pb.DeleteUserResponse{
		Message: "Kullanıcı başarıyla silindi",
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

	port := getEnv("PORT", ":50051")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("TCP dinlenemedi: %v", err)
	}
	log.Printf("Sunucu %s adresinde dinleniyor", port)

	jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

	// gRPC sunucusunu Prometheus için enstrümante et
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterUserServiceServer(s, &server{db: db, jwtSecret: jwtSecret})
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
		log.Fatalf("Sunucu hizmet vermedi: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
