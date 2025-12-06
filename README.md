# Atlas 2.0 - E-Ticaret Mikroservis Projesi

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)
![gRPC](https://img.shields.io/badge/gRPC-Protobuf-4285F4?style=flat&logo=google)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Helm-326CE5?style=flat&logo=kubernetes)

Atlas 2.0, eski monolitik e-ticaret sisteminin yerini almak üzere tasarlanmış; ölçeklenebilir, yüksek performanslı ve modern bir mikroservis mimarisidir. Bu proje **Go (Golang)**, **gRPC**, **PostgreSQL** ve **React** kullanılarak geliştirilmiştir.

## 📋 İçindekiler

1. [Proje Mimarisi](#-proje-mimarisi)
2. [Teknoloji Yığını](#-teknoloji-yığını)
3. [Servisler ve Özellikleri](#-servisler-ve-özellikleri)
4. [Veritabanı Yapısı](#-veritabanı-yapısı)
5. [Kurulum ve Çalıştırma](#-kurulum-ve-çalıştırma)
    - [Docker Compose ile Çalıştırma](#docker-compose-ile-çalıştırma-geliştirme-ortamı)
    - [Kubernetes (Helm) ile Dağıtım](#kubernetes-helm-ile-dağıtım)
6. [İzleme ve Metrikler (Monitoring)](#-i%CC%87zleme-ve-metrikler-monitoring)
7. [Proje Yapısı](#-proje-yapısı)  
8. [Güvenlik](#-güvenlik)
9. [Katkıda Bulunma](#%E2%80%8D-katk%C4%B1da-bulunma)

---

## 🏗 Proje Mimarisi

Sistem, API-First yaklaşımıyla tasarlanmış olup "Headless Commerce" yapısındadır. Frontend ve Backend tamamen ayrılmıştır. Servisler arası iletişimde performans için **gRPC**, dış dünya ile iletişimde ise **REST (HTTP/JSON)** kullanılmaktadır.

### Mimari Bileşenler:
* **Frontend (React):** Kullanıcı, Satıcı ve Yönetici panellerini içeren arayüz.
* **API Gateway:** HTTP isteklerini karşılar, JWT doğrulaması yapar ve ilgili mikroservise gRPC üzerinden yönlendirir.
* **Microservices:** User, Product ve Order servisleri kendi veritabanlarına sahiptir.

> **Not:** Mimari diyagramlar `projectFiles/` klasörü altında `.drawio` formatında mevcuttur.

---

## 🛠 Teknoloji Yığını

* **Backend:** Go (Golang) 1.22
* **İletişim Protokolleri:** gRPC (Protobuf), HTTP/1.1 (REST)
* **Veritabanı:** PostgreSQL 15 (Her servis için izole veritabanı)
* **Frontend:** React.js, Axios, Context API
* **Altyapı & DevOps:** Docker, Docker Compose, Kubernetes, Helm Charts, GitHub Actions
* **Monitoring:** Prometheus, Grafana (Özel Dashboardlar dahil)

---

## 📦 Servisler ve Özellikleri

| Servis | Port (HTTP/gRPC) | Açıklama | Veritabanı |
| :--- | :--- | :--- | :--- |
| **Frontend** | 3000 (Web) | React tabanlı kullanıcı arayüzü. Rol bazlı (Admin/User/Seller) sayfa yönetimi. | - |
| **Gateway** | 8080 (HTTP) | Merkezi giriş noktası. Authentication (JWT), Routing ve Protocol Translation (HTTP -> gRPC). | - |
| **User Service** | 50051 (gRPC) | Kayıt, Giriş, Profil Yönetimi, Rol Yönetimi (Kullanıcı, Satıcı, Yönetici). | `userdb` |
| **Product Service** | 50052 (gRPC) | Ürün listeleme, ekleme, güncelleme, stok takibi. | `productdb` |
| **Order Service** | 50053 (gRPC) | Sipariş oluşturma, sipariş listeleme, kargo durumu güncelleme. | `orderdb` |

---

## 🗄 Veritabanı Yapısı

Her mikroservis, veri tutarlılığı ve bağımsızlık ilkesi gereği kendi PostgreSQL veritabanına sahiptir.

1.  **User DB (`userdb`):**
    * `User`: Kullanıcı bilgileri, şifre (hash), rol.
    * `Address`: Kullanıcı adresleri.
2.  **Product DB (`productdb`):**
    * `Product`: Ürün detayları, fiyat, stok, satıcı ID'si.
    * `Category`: Ürün kategorileri.
    * `Currency`: Para birimleri.
3.  **Order DB (`orderdb`):**
    * `Order`: Sipariş başlık bilgileri, toplam tutar, durum.
    * `OrderDetail`: Sipariş kalemleri, miktar, birim fiyat.

---

## 🚀 Kurulum ve Çalıştırma

### Ön Koşullar
* Docker ve Docker Compose
* (Opsiyonel) Kubernetes Cluster (Minikube veya Kind) ve Helm

### Docker Compose ile Çalıştırma (Geliştirme Ortamı)

En hızlı kurulum yöntemidir. Tüm servisleri ve veritabanlarını ayağa kaldırır.

1.  Repoyu klonlayın:
    ```bash
    git clone [https://github.com/deverenozcan/atlas2.0.git](https://github.com/deverenozcan/atlas2.0.git)
    cd atlas2.0
    ```

2.  Docker Compose'u başlatın:
    ```bash
    docker-compose up --build
    ```

3.  Servislerin açılmasını bekleyin (Healthcheckler tanımlıdır). Erişim adresleri:
    * **Frontend:** [http://localhost:3000](http://localhost:3000)
    * **Gateway API:** [http://localhost:8080/api](http://localhost:8080/api)
    * **Prometheus Metrics:** [http://localhost:9090/metrics](http://localhost:9090/metrics) (Port yönlendirmesine göre değişebilir)

### Kubernetes (Helm) ile Dağıtım

Proje, `atlas2-helm-chart` klasöründe Helm chartlarını barındırır.

1.  Helm chart'ı yükleyin:
    ```bash
    helm install atlas2 ./atlas2-helm-chart
    ```

2.  Ingress ayarlarını `values.yaml` üzerinden yapılandırarak `atlas2.deverenozcan.com` (veya kendi hostunuz) üzerinden erişebilirsiniz.

---

## 📊 İzleme ve Metrikler (Monitoring)

Proje, **Prometheus** metriklerini dışa açacak şekilde yapılandırılmıştır (`github.com/grpc-ecosystem/go-grpc-prometheus`).

* **ServiceMonitor:** Kubernetes üzerinde Prometheus Operator ile otomatik keşif için `ServiceMonitor` tanımları mevcuttur.
* **Grafana Dashboard:** `atlas2-helm-chart/files` altında hazır dashboard JSON dosyaları bulunur:
    * `atlas2-dashboard.json`: Genel mikroservis metrikleri.
    * `todolist-dashboard.json`: Örnek dashboard.

---

## 📂 Proje Yapısı

````

Atlas2.0/
├── .github/workflows/   \# CI/CD Pipeline (Docker Build & Push)
├── atlas2-helm-chart/   \# Kubernetes Helm Chart dosyaları
│   ├── files/           \# DB init scriptleri ve Dashboard JSON'ları
│   └── templates/       \# K8s Deployment, Service, Ingress, ConfigMap vb.
├── docker-compose.yml   \# Yerel geliştirme ortamı konfigürasyonu
├── frontend/            \# React Uygulaması
│   ├── public/
│   └── src/
│       ├── components/  \# Login, ProductList vb. bileşenler
│       ├── context/     \# AuthContext, CartContext
│       ├── pages/       \# Dashboard, Orders vb. sayfalar
│       └── services/    \# Axios API çağrıları
├── gateway/             \# API Gateway (Go)
│   └── proto/           \# Ortak Proto dosyaları
├── order/               \# Order Microservice (Go)
├── product/             \# Product Microservice (Go)
├── user/                \# User Microservice (Go)
└── projectFiles/        \# Mimari çizimler ve gereksinim dokümanları

````

---

## 🔒 Güvenlik

* **JWT (JSON Web Token):** Kullanıcı oturum yönetimi için kullanılır.
* **Password Hashing:** Kullanıcı şifreleri `bcrypt` ile hashlenerek saklanır.
* **Role Based Access Control (RBAC):** Frontend ve Gateway seviyesinde rol kontrolü (Admin, Satıcı, Kullanıcı) yapılır.

---

## 👨‍💻 Katkıda Bulunma

1.  Forklayın.
2.  Feature branch oluşturun (`git checkout -b feature/yeni-ozellik`).
3.  Değişikliklerinizi commit edin (`git commit -m 'Yeni özellik eklendi'`).
4.  Branch'inizi pushlayın (`git push origin feature/yeni-ozellik`).
5.  Pull Request oluşturun.

---
*Bu proje [deverenozcan](https://github.com/deverenozcan) tarafından geliştirilmiştir.*
```
