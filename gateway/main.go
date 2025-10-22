package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "gateway/proto/order"
	productpb "gateway/proto/product"
	userpb "gateway/proto/user"
)

func main() {
	// gRPC bağlantıları (uygulama açılışında oluşturulur)
	userConn, err := grpc.NewClient(getEnv("USER_SERVICE_ADDR", "user:50051"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("user service'e bağlanılamadıı.: %v", err)
	}
	defer func() {
		if cerr := userConn.Close(); cerr != nil {
			log.Printf("userConn kapatma hatası: %v", cerr)
		}
	}()
	userClient := userpb.NewUserServiceClient(userConn)

	productConn, err := grpc.NewClient(getEnv("PRODUCT_SERVICE_ADDR", "product:50052"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("product service'e bağlanılamadı: %v", err)
	}
	defer func() {
		if cerr := productConn.Close(); cerr != nil {
			log.Printf("productConn kapatma hatası: %v", cerr)
		}
	}()
	productClient := productpb.NewProductServiceClient(productConn)

	orderConn, err := grpc.NewClient(getEnv("ORDER_SERVICE_ADDR", "order:50053"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("order service'e bağlanılamadı: %v", err)
	}
	defer func() {
		if cerr := orderConn.Close(); cerr != nil {
			log.Printf("orderConn kapatma hatası: %v", cerr)
		}
	}()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handleListUsers(w, r, userClient)
	}).Methods("GET", "OPTIONS")

	// Product routes
	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		handleListProducts(w, r, productClient)
	}).Methods("GET", "OPTIONS")

	// Order routes
	r.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		handleListOrders(w, r, orderClient)
	}).Methods("GET", "OPTIONS")

	handler := corsMiddleware(r)

	log.Println("Gateway server 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleListUsers(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.ListUsers(ctx, &userpb.ListUsersRequest{})
	if err != nil {
		http.Error(w, "Kullanıcıları listelerken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.GetUsers()); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleListProducts(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		http.Error(w, "Ürünleri listelerken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.GetProducts()); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleListOrders(w http.ResponseWriter, r *http.Request, client orderpb.OrderServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.ListOrders(ctx, &orderpb.ListOrdersRequest{})
	if err != nil {
		http.Error(w, "Siparişleri listelerken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.GetOrders()); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Basit CORS middleware'i
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
