package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "gateway/proto/order"
	productpb "gateway/proto/product"
	userpb "gateway/proto/user"
)

func main() {
	userConn, err := grpc.NewClient(getEnv("USER_SERVICE_ADDR", "user:50051"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("user service'e bağlanilamadi: %v", err)
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
	// Basic API info endpoint for both /api and /api/
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"api": "ok"})
	}).Methods("GET", "OPTIONS")

	// Serve all gateway endpoints under /api so Ingress that forwards /api will match routes (e.g. /api/auth/login)
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"api": "ok"})
	}).Methods("GET", "OPTIONS")

	api.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, userClient)
	}).Methods("POST", "OPTIONS")

	api.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, userClient)
	}).Methods("POST", "OPTIONS")

	api.HandleFunc("/users", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleListUsers(w, r, userClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/users/{id}", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleGetUser(w, r, userClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/users/{id}", adminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleUpdateUser(w, r, userClient)
	})).Methods("PUT", "OPTIONS")

	api.HandleFunc("/users/{id}", adminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleDeleteUser(w, r, userClient)
	})).Methods("DELETE", "OPTIONS")

	api.HandleFunc("/products", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleListProducts(w, r, productClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/products/my", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleListMyProducts(w, r, productClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/products/{id}", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleGetProduct(w, r, productClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/products", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleCreateProduct(w, r, productClient)
	})).Methods("POST", "OPTIONS")

	api.HandleFunc("/products/{id}", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleUpdateProduct(w, r, productClient)
	})).Methods("PUT", "OPTIONS")

	api.HandleFunc("/products/{id}", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleDeleteProduct(w, r, productClient)
	})).Methods("DELETE", "OPTIONS")

	api.HandleFunc("/orders", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleListOrders(w, r, orderClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/orders/my", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleListMyOrders(w, r, orderClient)
	})).Methods("GET", "OPTIONS")

	api.HandleFunc("/orders", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleCreateOrder(w, r, orderClient)
	})).Methods("POST", "OPTIONS")

	api.HandleFunc("/orders/{id}/status", adminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handleUpdateOrderStatus(w, r, orderClient)
	})).Methods("PATCH", "OPTIONS")

	handler := corsMiddleware(r)

	log.Println("Gateway server 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleRegister(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Ad       string `json:"ad"`
		Soyad    string `json:"soyad"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.Register(ctx, &userpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Ad:       req.Ad,
		Soyad:    req.Soyad,
	})
	if err != nil {
		http.Error(w, "Kayıt başarısız: "+err.Error(), http.StatusInternalServerError)
		return
	}

	registerResponse := map[string]interface{}{
		"user_id": res.UserId,
		"message": res.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(registerResponse); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.Login(ctx, &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "Giriş başarısız: "+err.Error(), http.StatusUnauthorized)
		return
	}

	loginResponse := map[string]interface{}{
		"token":   res.Token,
		"user_id": res.UserId,
		"email":   res.Email,
		"role":    res.Role.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
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

	users := make([]map[string]interface{}, len(res.GetUsers()))
	for i, user := range res.GetUsers() {
		users[i] = map[string]interface{}{
			"user_id":    user.UserId,
			"email":      user.Email,
			"ad":         user.Ad,
			"soyad":      user.Soyad,
			"role":       user.Role.String(),
			"address_id": user.AddressId,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID'si", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &userpb.GetUserRequest{UserId: id})
	if err != nil {
		http.Error(w, "Kullanıcı getirilemedi: "+err.Error(), http.StatusNotFound)
		return
	}

	user := res.GetUser()
	userJSON := map[string]interface{}{
		"user_id":    user.UserId,
		"email":      user.Email,
		"ad":         user.Ad,
		"soyad":      user.Soyad,
		"role":       user.Role.String(),
		"address_id": user.AddressId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userJSON); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID'si", http.StatusBadRequest)
		return
	}

	var req struct {
		Email string `json:"email"`
		Ad    string `json:"ad"`
		Soyad string `json:"soyad"`
		Role  string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	var roleEnum userpb.UserRole
	switch req.Role {
	case "KULLANICI":
		roleEnum = userpb.UserRole_KULLANICI
	case "SATICI":
		roleEnum = userpb.UserRole_SATICI
	case "YONETICI":
		roleEnum = userpb.UserRole_YONETICI
	default:
		http.Error(w, "Geçersiz rol değeri", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.UpdateUser(ctx, &userpb.UpdateUserRequest{
		UserId: id,
		Email:  req.Email,
		Ad:     req.Ad,
		Soyad:  req.Soyad,
		Role:   roleEnum,
	})
	if err != nil {
		http.Error(w, "Kullanıcı güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, client userpb.UserServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID'si", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.DeleteUser(ctx, &userpb.DeleteUserRequest{UserId: id})
	if err != nil {
		http.Error(w, "Kullanıcı silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
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

func handleListMyProducts(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]
	jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token okunamadı", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Claims okunamadı", http.StatusUnauthorized)
		return
	}

	userID, _ := claims["user_id"].(float64)
	role, _ := claims["role"].(string)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		http.Error(w, "Ürünleri listelerken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	products := res.GetProducts()

	if role != "YONETICI" {
		filteredProducts := make([]*productpb.Product, 0)
		for _, product := range products {
			if product.SellerId == int64(userID) {
				filteredProducts = append(filteredProducts, product)
			}
		}
		products = filteredProducts
	}

	if products == nil {
		products = make([]*productpb.Product, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetProduct(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz ürün ID'si", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.GetProduct(ctx, &productpb.GetProductRequest{ProductId: id})
	if err != nil {
		http.Error(w, "Ürün getirilemedi: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.GetProduct()); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req struct {
		ProductName   string  `json:"product_name"`
		Description   string  `json:"description"`
		Stock         int32   `json:"stock"`
		Price         float64 `json:"price"`
		Currency      string  `json:"currency"`
		SellerID      int64   `json:"seller_id"`
		CategoryID    int64   `json:"category_id"`
		DimensDetails string  `json:"dimens_details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.CreateProduct(ctx, &productpb.CreateProductRequest{
		ProductName:   req.ProductName,
		Description:   req.Description,
		Stock:         req.Stock,
		Price:         req.Price,
		Currency:      req.Currency,
		SellerId:      req.SellerID,
		CategoryId:    req.CategoryID,
		DimensDetails: req.DimensDetails,
	})
	if err != nil {
		http.Error(w, "Ürün oluşturulamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz ürün ID'si", http.StatusBadRequest)
		return
	}

	var req struct {
		ProductName   string  `json:"product_name"`
		Description   string  `json:"description"`
		Stock         int32   `json:"stock"`
		Price         float64 `json:"price"`
		Currency      string  `json:"currency"`
		CategoryID    int64   `json:"category_id"`
		DimensDetails string  `json:"dimens_details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.UpdateProduct(ctx, &productpb.UpdateProductRequest{
		ProductId:     id,
		ProductName:   req.ProductName,
		Description:   req.Description,
		Stock:         req.Stock,
		Price:         req.Price,
		Currency:      req.Currency,
		CategoryId:    req.CategoryID,
		DimensDetails: req.DimensDetails,
	})
	if err != nil {
		http.Error(w, "Ürün güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request, client productpb.ProductServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz ürün ID'si", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.DeleteProduct(ctx, &productpb.DeleteProductRequest{ProductId: id})
	if err != nil {
		http.Error(w, "Ürün silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request, client orderpb.OrderServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]
	jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token okunamadı", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Claims okunamadı", http.StatusUnauthorized)
		return
	}

	userID, _ := claims["user_id"].(float64)

	var req struct {
		Items []struct {
			ProductID int64   `json:"product_id"`
			Quantity  int32   `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	var orderItems []*orderpb.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, &orderpb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId: int64(userID),
		Items:  orderItems,
	})
	if err != nil {
		http.Error(w, "Sipariş oluşturulamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
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

func handleListMyOrders(w http.ResponseWriter, r *http.Request, client orderpb.OrderServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]
	jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token okunamadı", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Claims okunamadı", http.StatusUnauthorized)
		return
	}

	userID, _ := claims["user_id"].(float64)
	role, _ := claims["role"].(string)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.ListOrders(ctx, &orderpb.ListOrdersRequest{})
	if err != nil {
		http.Error(w, "Siparişleri listelerken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	orders := res.GetOrders()

	if role != "YONETICI" {
		filteredOrders := make([]*orderpb.Order, 0)
		for _, order := range orders {
			if order.UserId == int64(userID) {
				filteredOrders = append(filteredOrders, order)
			}
		}
		orders = filteredOrders
	}

	if orders == nil {
		orders = make([]*orderpb.Order, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdateOrderStatus(w http.ResponseWriter, r *http.Request, client orderpb.OrderServiceClient) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Geçersiz sipariş ID'si", http.StatusBadRequest)
		return
	}

	var req struct {
		IsShipped bool `json:"is_shipped"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res, err := client.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		OrderId:   id,
		IsShipped: req.IsShipped,
	})
	if err != nil {
		http.Error(w, "Sipariş durumu güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSON kodlama hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// JWT Middleware - Token doğrulama
func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header gerekli", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Geçersiz Authorization formatı", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz veya süresi dolmuş token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header gerekli", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Geçersiz Authorization formatı", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		jwtSecret := getEnv("JWT_SECRET", "default-secret-key-change-in-production")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz veya süresi dolmuş token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token claims okunamadı", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Token'da role bilgisi bulunamadı", http.StatusUnauthorized)
			return
		}

		if role != "YONETICI" {
			http.Error(w, "Bu işlem için yönetici yetkisi gereklidir", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
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
