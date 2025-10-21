### 1. GİRİŞ VE KAPSAM

### 1.1. Proje Amacı

Mevcut "Atlas" e-ticaret sistemi, monolitik ve eskiyen mimarisi nedeniyle artan kullanıcı trafiğini ve eş zamanlı işlem yükünü karşılayamamaktadır. Bu projenin temel amacı, "Atlas" sistemini, yüksek ölçeklenebilirlik, performans, güvenlik ve sürdürülebilirlik ilkelerine dayanan modern bir mimari (Mikroservis, API-First, Headless) kullanarak sıfırdan inşa etmektir. Yeni sistem, mevcut iş fonksiyonlarını desteklemenin yanı sıra gelecekteki büyüme ve yeni pazar gereksinimlerine hızla adapte olabilecek esnek bir temel sağlayacaktır.

### 1.2. İş Hedefleri

- **Ölçeklenebilirlik:** Mevcut kullanıcı trafiğinin en az 10 katını sorunsuz bir şekilde yönetebilmek. (Örn: Kara Cuma gibi kampanya dönemlerinde 100.000 eş zamanlı kullanıcıyı desteklemek).
- **Performans:** Ortalama sayfa yükleme süresini < 2 saniye, API yanıt sürelerini (p95) < 200ms seviyesine indirmek.
- **Çeviklik (Agility):** Yeni özellikleri ve modülleri, sistemin tamamını durdurmadan (zero-downtime) ve hızlı bir şekilde (günlük/haftalık sürümlerle) devreye alabilmek.
- **Güvenilirlik:** %99.95 oranında çalışma süresi (uptime) sağlamak ve sistemin bir bileşenindeki hatanın diğer bileşenleri etkilemesini (cascading failure) önlemek.
- **Maliyet Optimizasyonu:** Kaynak kullanımını optimize ederek (örn: talep bazlı otomatik ölçeklenme) altyapı maliyetlerini düşürmek.

### 1.3. Kapsam

**Kapsam Dahilinde (In-Scope):**

- Tüm temel e-ticaret işlevlerinin (Kullanıcı, Ürün, Sepet, Sipariş, Ödeme) yeniden tasarlanması ve geliştirilmesi.
- Mikroservis tabanlı, bulut-yerel (cloud-native) bir arka uç mimarisinin kurulması.
- Tüm işlevselliği dışa açan bir API ağ geçidinin (API Gateway) tasarlanması.
- Modern, duyarlı (responsive) bir ön uç (Web) ve/veya mobil uygulamaları (iOS/Android) besleyecek "Headless" (Bağımsız Arayüz) bir altyapı.
- Kritik verilerin (Kullanıcılar, Ürünler, Siparişler) eski sistemden yeni sisteme taşınması (Data Migration).
- Üçüncü parti sistemlerle (Ödeme, Kargo, ERP) entegrasyonların modern standartlarda (API, Webhook) yeniden yapılması.

**Kapsam Dışında (Out-of-Scope):**

- Mevcut monolitik sistemin "yamalarla" iyileştirilmesi (Proje hedefi "yeniden inşa"dır).
- İç Lojistik ve Depo Yönetim Sistemi (WMS) yazılımının geliştirilmesi (Ancak entegrasyonu kapsam dahilindedir).
- 10 yıllık arşivsel sipariş verilerinin tamamının yeni operasyonel veritabanına taşınması (Bu veriler bir "data warehouse"a aktarılabilir).

### 1.4. Paydaşlar

- **Son Kullanıcılar:** (Müşteriler - B2C/B2B)
- **İş Birimleri:** Kategori Yönetimi, Pazarlama, Operasyon (Sipariş Yönetimi), Müşteri Hizmetleri
- **Teknik Ekip:** Yazılım Geliştirme, DevOps, Veritabanı Yöneticileri
- **Yönetim:** CEO, CTO, CFO
- **Dış Ortaklar:** Ödeme Sağlayıcılar, Lojistik Firmaları, Tedarikçiler

---

### 2. SİSTEM MİMARİSİ VE TASARIM YAKLAŞIMI

Bu bölüm, projenin "modern yaklaşım" temelini oluşturur.

### 2.1. Mimari Yaklaşım

- **Mikroservis Mimarisi:** Sistem, her biri belirli bir iş alanından (domain) sorumlu, bağımsız olarak geliştirilebilen, dağıtılabilen ve ölçeklenebilen küçük servislere bölünecektir.
- **API-First (API Öncelikli) Tasarım:** Tüm sistem işlevselliği, iç ve dış tüketiciler için (Web, Mobil, 3. Parti) güvenli, belgelenmiş ve standartlaştırılmış API'ler (RESTful veya gRPC) aracılığıyla sunulacaktır.
- **Headless Commerce (Bağımsız Arayüz):** Arka uç (iş mantığı, veritabanı) ile ön uç (kullanıcı arayüzü) tamamen ayrılacaktır. Bu, farklı kanallar (Web, Mobil Uygulama, Kiosk, IoT) için kolayca arayüzler geliştirmeyi sağlar.
- **Olay Güdümlü (Event-Driven) İletişim:** Servisler arası gevşek bağlılığı (loose coupling) sağlamak için, özellikle sipariş akışı gibi asenkron işlemler için bir mesajlaşma kuyruğu (örn: Kafka, RabbitMQ) kullanılacaktır. (Örn: "Sipariş Oluşturuldu" olayı yayınlanır, "Stok" ve "Bildirim" servisleri bu olayı dinler).

### 2.2. Teknoloji Ekosistemi (Önerilen)

- **Konteynerizasyon & Orkestrasyon:** Docker & Kubernetes (K8s) (Bulut-yerel dağıtım ve otomatik ölçeklenme için).
- **Veri Depolama (Polyglot Persistence):** Her servisin kendi ihtiyacına uygun veritabanı teknolojisini kullanması.
    - *İlişkisel (SQL):* Sipariş, Kullanıcı servisleri (örn: PostgreSQL, MySQL).
    - *NoSQL (Belge):* Ürün kataloğu, İçerik servisi (örn: MongoDB, DynamoDB).
    - *NoSQL (Key-Value/Cache):* Sepet, Oturum (Session) servisi (örn: Redis, Memcached).
    - *Arama Motoru:* Ürün arama servisi (örn: Elasticsearch, Algolia).
- **API Gateway:** Güvenlik (Authentication, Rate Limiting), yönlendirme ve API birleştirme işlemleri için (örn: Kong, AWS API Gateway, Nginx).
- **DevOps & CI/CD:** GitLab CI, Jenkins, GitHub Actions kullanılarak tam otomatik test ve dağıtım (deployment) süreçleri.
- **İzleme (Monitoring) & Günlükleme (Logging):** Prometheus, Grafana (Metrikler) ve ELK Stack/Datadog (Merkezi Loglama).

---

### 3. FONKSİYONEL GEREKSİNİMLER (Functional Requirements)

Sistem, aşağıdaki mikroservisler (veya iş alanları) etrafında gruplanan işlevleri yerine getirmelidir.

### 3.1. Kullanıcı & Kimlik Servisi (Identity Service)

- **FR-3.1.1:** Kullanıcılar e-posta/şifre ile veya sosyal medya hesapları (Google, Apple) ile kayıt olabilmelidir (OAuth 2.0).
- **FR-3.1.2:** Kullanıcılar sisteme giriş yapabilmeli ve JWT (JSON Web Token) tabanlı oturum yönetimi sağlanmalıdır.
- **FR-3.1.3:** Kullanıcılar "Şifremi Unuttum" akışını tamamlayabilmelidir.
- **FR-3.1.4:** Kullanıcılar profil bilgilerini (Ad, Soyad, Telefon) yönetebilmelidir.
- **FR-3.1.5:** Kullanıcılar birden fazla teslimat ve fatura adresi tanımlayabilmelidir.
- **FR-3.1.6:** Rol bazlı yetkilendirme (RBAC) altyapısı bulunmalıdır (örn: User, Admin, ProductOwner).

### 3.2. Ürün Katalog Servisi (Catalog Service)

- **FR-3.2.1:** ProductOwner yeni ürünler, kategoriler ve markalar oluşturabilmeli (CRUD).
- **FR-3.2.2:** Ürünler birden fazla resim, detaylı açıklama, teknik özellik (nitelik) ve varyant (Renk, Beden) bilgisi içerebilmelidir.
- **FR-3.2.3:** Kullanıcılar ürünleri listeleyebilmelidir (Kategori, Marka bazlı).
- **FR-3.2.4:** Kullanıcılar ürünler üzerinde detaylı filtreleme (Fiyat aralığı, Renk, Marka vb.) ve sıralama yapabilmelidir.
- **FR-3.2.5:** Kullanıcılar ürünleri anahtar kelime ile arayabilmelidir (Otomatik tamamlama destekli).

### 3.3. Stok Yönetim Servisi (Inventory Service)

- **FR-3.3.1:** Her ürünün/varyantın stok adedi gerçek zamanlı olarak tutulmalıdır.
- **FR-3.3.2:** Ödeme işlemi tamamlandığında stok adedi otomatik olarak düşülmelidir.
- **FR-3.3.3:** İptal/iade durumlarında stok otomatik olarak güncellenmelidir.

### 3.4. Sepet Servisi (Cart/Basket Service)

- **FR-3.4.1:** Kullanıcılar sepete ürün ekleyebilmeli, çıkarabilmeli ve adedini güncelleyebilmelidir.

### 3.5. Sipariş ve Ödeme Servisi (Order & Payment Service)

- **FR-3.5.1:** Kullanıcılar sepetlerini, adres ve ödeme adımlarını içeren bir "checkout" süreci ile siparişe dönüştürebilmelidir.
- **FR-3.5.5:** Başarılı ödeme sonrası sipariş "Onaylandı" statüsüne geçmeli ve benzersiz bir sipariş numarası üretilmelidir.
- **FR-3.5.6:** Kullanıcılar "Siparişlerim" sayfasından siparişlerinin geçmişini ve güncel durumunu (Alındı, Hazırlanıyor, Kargoda, Teslim Edildi, İptal) takip edebilmelidir.
- **FR-3.5.7:** Sipariş süreci (Ödeme, Stok düşme) "Saga" deseni gibi hata telafi mekanizmalarıyla (compensation) yönetilmelidir. (Örn: Ödeme alındı ama stok düşülemediyse, ödeme iade edilmeli).
- **FR-3.5.8:** Sipariş iptal ve iade talepleri sistem üzerinden yönetilebilmelidir.

### 3.6. Bildirim Servisi (Notification Service)

- **FR-3.6.1:** Sistem, kritik olaylarda (Sipariş Onayı, Kargoya Verildi, Şifre Sıfırlama) kullanıcılara e-posta göndermelidir.

### 3.7. Yönetim Paneli (Admin Panel)

- **FR-3.7.1:** Yetkili kullanıcılar (Adminler) için ayrı bir yönetim arayüzü olmalıdır.
- **FR-3.7.2:** Panel üzerinden Ürün, Kategori, Kullanıcı, Sipariş yönetimi (CRUD) yapılabilmelidir.
- **FR-3.7.3:** Siparişlerin durumu (örn: "Hazırlanıyor" -> "Kargoya Verildi") manuel olarak güncellenebilmeli ve kargo takip numarası girilebilmelidir.

---

### 4. FONKSİYONEL OLMAYAN GEREKSİNİMLER (Non-Functional Requirements - NFRs)

Bu bölüm, sistemin "nasıl" çalışması gerektiğini tanımlar ve monolitik yapıdaki temel sorunları (performans, ölçeklenebilirlik) doğrudan hedefler.

### 4.1. Performans

- **NFR-4.1.1 (Yanıt Süresi):** API ağ geçidinden ölçülen p95 (kullanıcıların %95'inin deneyimlediği) API yanıt süresi < 200ms olmalıdır (Arama ve ödeme API'leri hariç).
- **NFR-4.1.2 (Sayfa Yükleme):** Google Core Web Vitals (LCP, FID, CLS) "İyi" seviyede olmalıdır. Kritik sayfalar (Ana sayfa, Kategori, Ürün Detay) 2 saniye altında tam olarak yüklenebilir olmalıdır.
- **NFR-4.1.3 (Yük Altında):** Sistem, 1.000 (bin) anlık sipariş (checkout) işlemini 1 dakika içinde sorunsuz tamamlayabilmelidir.

### 4.2. Ölçeklenebilirlik

- **NFR-4.2.1 (Yatay Ölçeklenme):** Tüm mikroservisler yatay olarak (yeni pod'lar/instancelar eklenerek) ölçeklenebilmelidir.
- **NFR-4.2.2 (Otomatik Ölçeklenme):** Kubernetes (HPA) veya bulut sağlayıcının (AWS Auto Scaling) özellikleri kullanılarak, CPU/RAM veya trafik yüküne göre servisler otomatik olarak ölçeklenmeli ve talep düşünce küçülmelidir (Maliyet optimizasyonu).
- **NFR-4.2.3 (Veritabanı):** Veritabanları, okuma replikaları (read replicas) ile ölçeklenebilmeli ve yoğun yazma yükünü kaldırabilmelidir.

### 4.3. Güvenilirlik ve Erişilebilirlik (Reliability & Availability)

- **NFR-4.3.1 (Çalışma Süresi):** Sistemin genel çalışma süresi %99.95 (yıllık ~4.38 saat kesinti) olmalıdır.
- **NFR-4.3.2 (Dayanıklılık - Resilience):** Bir servisin (örn: "Yorum Servisi") çökmesi, sistemin ana işlevini (örn: "Sipariş Verme") engellememelidir (Circuit Breaker, Bulkhead desenleri kullanılmalıdır).
- **NFR-4.3.3 (Felaket Kurtarma):** Veritabanları düzenli olarak yedeklenmeli ve farklı bir coğrafi bölgede (Multi-Region/AZ) kurtarma planı (Disaster Recovery Plan) bulunmalıdır.

### 4.4. Güvenlik

- **NFR-4.4.1:** OWASP Top 10 zafiyetlerine (XSS, SQL Injection, CSRF vb.) karşı koruma sağlanmalıdır.
- **NFR-4.4.2:** Tüm iletişim (kullanıcı-sunucu ve servisler arası) SSL/TLS (HTTPS) ile şifrelenmelidir.
- **NFR-4.4.3:** Kullanıcı şifreleri, güçlü ve tuzlanmış (salted) hash algoritmaları (örn: bcrypt, Argon2) ile saklanmalıdır.
- **NFR-4.4.4:** API'ler, API Gateway üzerinden kimlik doğrulama (Authentication) ve yetkilendirme (Authorization) mekanizmalarından geçmelidir.
- **NFR-4.4.5:** API'ler için hız limiti (Rate Limiting) ve DDoS koruması uygulanmalıdır.

### 4.5. Sürdürülebilirlik (Maintainability)

- **NFR-4.5.1:** Her mikroservisin kendi bağımsız CI/CD pipeline'ı olmalı ve diğer servislerden bağımsız deploy edilebilmelidir.
- **NFR-4.5.2:** Kod kalitesi için statik kod analizi ve %80 üzeri birim test (unit test) kapsamı hedeflenmelidir.
- **NFR-4.5.3:** Tüm servisler, merkezi bir loglama sistemine (örn: ELK Stack) yapılandırılmış (structured) loglar göndermelidir.
- **NFR-4.5.4:** Tüm servislerin sağlık durumu (health check) ve temel metrikleri (CPU, RAM, İstek Sayısı, Hata Oranı) merkezi bir izleme sistemi (örn: Prometheus/Grafana) üzerinden takip edilmelidir.

---

---