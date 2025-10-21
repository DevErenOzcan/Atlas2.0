### **Use Case 01: E-posta ile Kullanıcı Kaydı**

- **ID:** UC-001
- **Adı:** E-posta ile Yeni Kullanıcı Kaydı
- **Aktör:** Müşteri (Giriş yapmamış)
- **Ön Koşul:** Müşteri, "Kayıt Ol" sayfasındadır.
- **Temel Akış:**
    1. Müşteri ad, soyad, e-posta ve parola bilgilerini girer.
    2. Müşteri, "Kayıt Ol" butonuna tıklar.
    3. Sistem (Kullanıcı Servisi), e-postanın sistemde kayıtlı olup olmadığını kontrol eder.
    4. Sistem, parolayı güvenli bir şekilde hash'ler (örn: bcrypt).
    5. Sistem, yeni kullanıcıyı veritabanına kaydeder.
    6. Sistem, JWT (JSON Web Token) üretir ve kullanıcıyı oturum açmış duruma getirir.
    7. Sistem, kullanıcıyı ana sayfaya yönlendirir.
- **Alternatif Akışlar:**
    - **5a. E-posta zaten kayıtlı:** Sistem "Bu e-posta adresi zaten kullanılıyor." hatası gösterir. Kayıt işlemi durdurulur.
    - **5b. Hatalı/Eksik Bilgi:** Müşteri zorunlu bir alanı eksik doldurduysa, sistem uyarı mesajı gösterir.
- **Son Koşul:** Müşteri, sistemde yeni bir hesap oluşturmuş ve giriş yapmıştır.
- **İlgili Gereksinim:** FR-3.1.1, FR-3.1.2, FR-4.4.3

---

### **Use Case 02: Kayıtlı Kullanıcının E-posta/Şifre ile Giriş Yapması**

- **ID:** UC-002
- **Adı:** Kayıtlı Kullanıcının E-posta/Şifre ile Giriş Yapması
- **Aktör:** Müşteri
- **Ön Koşul:** Müşteri, "Giriş Yap" sayfasındadır.
- **Temel Akış:**
    1. Müşteri e-posta ve parola bilgilerini girer.
    2. Müşteri, "Giriş Yap" butonuna tıklar.
    3. Sistem (Kullanıcı Servisi), sağlanan parola ile veritabanındaki hashlenmiş parolayı karşılaştırır.
    4. Parolaların eşleştiğini doğrular.
    5. Sistem, bir JWT (JSON Web Token) oluşturur.
    6. Sistem, kullanıcıyı ana sayfaya yönlendirir.
- **Alternatif Akışlar:**
    - **3a. Hatalı Kimlik Bilgileri:** E-posta veya parola yanlışsa, sistem "E-posta veya parola hatalı" mesajını gösterir. İşlem durur.
- **Son Koşul:** Müşteri, kimlik doğrulaması (authenticated) yapılmış olarak sisteme erişebilir ve kişisel bilgilerine erişebilir.
- **İlgili Gereksinim:** FR-3.1.2, FR-3.1.4, FR-4.4.3, FR-4.4.4

---

### **Use Case 03: Kullanıcının Şifresini Sıfırlama Talebi**

- **ID:** UC-003
- **Adı:** Kullanıcının Şifresini Sıfırlama Talebi
- **Aktör:** Müşteri
- **Ön Koşul:** Müşteri, "Şifremi Unuttum" sayfasındadır.
- **Temel Akış:**
    1. Müşteri, kayıtlı e-posta adresini girer.
    2. Müşteri, "Gönder" butonuna tıklar.
    3. Sistem (Kullanıcı Servisi), e-posta adresinin kayıtlı olduğunu doğrular.
    4. Sistem, tek kullanımlık bir şifre sıfırlama token'ı (JWT veya benzersiz link) oluşturur.
    5. Sistem (Bildirim Servisi'ne asenkron çağrı), şifre sıfırlama linkini içeren bir e-posta gönderir.
    6. Müşteri e-postasındaki linke tıklar ve yeni şifresini belirler.
    7. Sistem, yeni şifreyi hashleyerek günceller.
- **Alternatif Akışlar:**
    - **3a. E-posta Kayıtlı Değil:** Sistem, güvenlik nedeniyle sadece "Eğer e-posta adresiniz sistemimizde kayıtlıysa, şifre sıfırlama e-postası gönderilecektir" şeklinde genel bir mesaj gösterir.
- **Son Koşul:** Kullanıcının parolası güvenli bir şekilde güncellenmiştir.
- **İlgili Gereksinim:** FR-3.1.3, FR-3.6.1, FR-4.4.3

---

### **Use Case 04: Ürün Detay Bilgilerini Görüntüleme**

- **ID:** UC-004
- **Adı:** Ürün Detay Bilgilerini Görüntüleme
- **Aktör:** Müşteri/Misafir
- **Ön Koşul:** Müşteri, bir ürün listeleme sayfasından veya aramadan bir ürün seçmiştir.
- **Temel Akış:**
    1. Müşteri ürün detay sayfasına tıklar.
    2. Sistem (API Gateway), Ürün Katalog Servisi'nden ürünün temel, teknik ve görsel bilgilerini ister.
    3. Katalog Servisi (NoSQL veritabanından), ürünün adını, fiyatını, detaylı açıklamasını, birden fazla resmini ve varsa varyantlarını (Renk, Beden) yanıtlar.
    4. Ön Uç, bu bilgileri hızlı bir şekilde (2 saniye altı) render eder.
- **Alternatif Akışlar:**
    - **2a. Ürün Bulunamadı:** Ürün ID'si geçersizse, Sistem "Ürün Bulunamadı (404)" hatası gösterir.
- **Son Koşul:** Müşteri, ürün hakkında tüm gerekli bilgilere erişmiş ve varyant seçimi yapabilir durumdadır.
- **İlgili Gereksinim:** FR-3.2.2, NFR-4.1.2

---

### **Use Case 05: Anahtar Kelime ile Ürün Arama**

- **ID:** UC-005
- **Adı:** Anahtar Kelime ile Ürün Arama
- **Aktör:** Müşteri/Misafir
- **Ön Koşul:** Müşteri, sitenin arama çubuğunu kullanmaktadır.
- **Temel Akış:**
    1. Müşteri, arama çubuğuna bir anahtar kelime yazmaya başlar (örn: "kırmızı elbise").
    2. Ön Uç, Arama Motoru'na (Elasticsearch/Algolia) otomatik tamamlama (autocomplete) isteği gönderir.
    3. Müşteri, arama kelimesini tamamlar ve "Ara" butonuna tıklar.
    4. Sistem (API Gateway $\to$ Arama Servisi), anahtar kelime ile arama sorgusunu çalıştırır.
    5. Arama Servisi, ilgili ürün listesini < 200ms içinde yanıtlar.
    6. Ön Uç, arama sonuçları sayfasını gösterir.
- **Alternatif Akışlar:**
    - **4a. Yavaş Arama:** Arama yanıt süresi 200ms üzerine çıkarsa, sistem performansı düşmüş kabul edilir.
- **Son Koşul:** Müşteri, arama kriterlerine uygun ürün listesine hızlıca ulaşmıştır.
- **İlgili Gereksinim:** FR-3.2.5, NFR-4.1.1

---

### **Use Case 06: Filtreleme ve Sıralama ile Ürün Listeleme**

- **ID:** UC-006
- **Adı:** Filtreleme ve Sıralama ile Ürün Listeleme
- **Aktör:** Müşteri/Misafir
- **Ön Koşul:** Müşteri, bir kategori sayfasındadır.
- **Temel Akış:**
    1. Müşteri, sol menüden bir filtre seçeneğine tıklar (örn: Fiyat Aralığı, Marka, Renk).
    2. Müşteri, üst menüden bir sıralama seçeneği (örn: En Düşük Fiyat) seçer.
    3. Ön Uç, yeni filtre ve sıralama parametrelerini içeren bir API isteğini Katalog Servisi'ne gönderir.
    4. Katalog Servisi, veritabanını (NoSQL) sorgular ve filtrelenmiş/sıralanmış ürün listesini döndürür.
- **Alternatif Akışlar:**
    - **4a. Filtreye Uyan Ürün Yok:** Seçilen kriterlere uyan ürün yoksa, sistem "Kriterlerinize uygun ürün bulunamadı" mesajı gösterir.
- **Son Koşul:** Müşteri, seçtiği kriterlere uygun olarak filtrelenmiş ve sıralanmış güncel ürün listesini görür.
- **İlgili Gereksinim:** FR-3.2.3, FR-3.2.4

---

### **Use Case 07: Ürünü Sepete Ekleme**

- **ID:** UC-007
- **Adı:** Ürünü Sepete Ekleme
- **Aktör:** Müşteri/Misafir
- **Ön Koşul:** Müşteri, bir ürün detay sayfasındadır ve varyant (varsa) seçimi yapmıştır.
- **Temel Akış:**
    1. Müşteri, "Sepete Ekle" butonuna tıklar.
    2. Ön Uç, ürün ID'si ve adedi ile Sepet Servisi'ne (Key-Value/Cache) bir istek gönderir.
    3. Sepet Servisi, ürünü kullanıcının sepetine ekler veya mevcut ürüne adedi ekler.
    4. Sepet Servisi, güncel sepet içeriğini ve toplam tutarı döndürür.
    5. Ön Uç, başarılı bir bildirim (örn: "Ürün sepete eklendi") gösterir.
- **Alternatif Akışlar:**
    - **3a. Stok Yetersizliği:** Stok Servisi'nden gelen bilgiye göre (asenkron veya senkron kontrol), eklenecek ürün adedi stoktan fazlaysa, Sepet Servisi işlemi reddeder ve "Maksimum stok adedine ulaşıldı" uyarısı gösterir.
- **Son Koşul:** Ürün, kullanıcının sepetine eklenmiştir.
- **İlgili Gereksinim:** FR-3.4.1, FR-3.3.1 (Stok kontrolü)

---

### **Use Case 08: Sepetteki Ürün Adedini Güncelleme**

- **ID:** UC-008
- **Adı:** Sepetteki Ürün Adedini Güncelleme
- **Aktör:** Müşteri/Misafir
- **Ön Koşul:** Müşteri, sepet sayfasındadır.
- **Temel Akış:**
    1. Müşteri, sepetindeki bir ürünün yanındaki adet alanını günceller (örn: 1'den 2 yapar).
    2. Ön Uç, Sepet Servisi'ne güncel ürün ID'si ve yeni adedi ile bir istek gönderir.
    3. Sepet Servisi, yeni adet için Stok Servisi'nden yeterli stok olup olmadığını kontrol eder.
    4. Stok yeterliyse, Sepet Servisi (Key-Value/Cache) sepeti günceller.
    5. Sepet Servisi, yeni toplam tutarı ve güncel sepet içeriğini döndürür.
- **Alternatif Akışlar:**
    - **3a. Stok Yetersizliği:** Yeni adet stoktan fazlaysa, Sepet Servisi işlemi reddeder, uyarı gösterir ve adedi maksimum stok sayısına ayarlar.
- **Son Koşul:** Sepetteki ürünün adedi güncellenmiş ve toplam tutar yeniden hesaplanmıştır.
- **İlgili Gereksinim:** FR-3.4.1, FR-3.3.1

---

### **Use Case 09: Çekiş (Checkout) Sürecini Başlatma**

- **ID:** UC-009
- **Adı:** Çekiş (Checkout) Sürecini Başlatma ve Adres Seçimi
- **Aktör:** Müşteri (Giriş Yapmış)
- **Ön Koşul:** Müşterinin sepeti dolu ve oturum açmıştır.
- **Temel Akış:**
    1. Müşteri, sepet sayfasında "Satın Al" butonuna tıklar.
    2. Sistem (API Gateway), Müşteriyi Teslimat ve Fatura Bilgileri sayfasına yönlendirir.
    3. Müşteri (Kullanıcı Servisi üzerinden), önceden tanımlı teslimat ve fatura adreslerinden birini seçer.
    4. Müşteri, kargo seçeneğini belirler ve "Ödeme Adımına Geç" butonuna tıklar.
    5. Sistem, seçilen adresleri Sipariş Servisi'ne geçirir.
- **Alternatif Akışlar:**
    - **3a. Adres Yok:** Müşteri önceden adres tanımlamadıysa, sistem yeni bir adres girmesini zorunlu kılar.
- **Son Koşul:** Teslimat ve fatura bilgileri seçilmiş/girilmiş, sipariş geçici olarak Sipariş Servisi'nde tutulmaya başlanmıştır.
- **İlgili Gereksinim:** FR-3.5.1, FR-3.1.5

---

### **Use Case 10: Başarılı Ödeme ile Sipariş Oluşturma**

- **ID:** UC-010
- **Adı:** Başarılı Ödeme ile Sipariş Oluşturma
- **Aktör:** Müşteri
- **Ön Koşul:** Müşteri, ödeme adımını tamamlamıştır.
- **Temel Akış:**
    1. Müşteri, ödeme bilgilerini girer ve "Siparişi Tamamla" butonuna tıklar.
    2. Sipariş Servisi, ödeme isteğini 3. parti Ödeme Sağlayıcısı'na gönderir.
    3. Ödeme başarılı olur. Sipariş Servisi, benzersiz bir Sipariş Numarası üretir.
    4. Sipariş Servisi, Siparişin durumunu "Ödeme Başarılı/Onaylandı" olarak günceller.
    5. Sipariş Servisi, "Sipariş Oluşturuldu" olayını mesajlaşma kuyruğuna (Event-Driven) yayınlar.
    6. Stok Servisi, olayı dinler ve ilgili ürünlerin stok adedini otomatik olarak düşer.
    7. Bildirim Servisi, olayı dinler ve müşteriye "Sipariş Onayı" e-postası gönderir.
    8. Ön Uç, müşteriyi "Sipariş Onaylandı" sayfasına yönlendirir.
- **Alternatif Akışlar:** Yok (Hata Akışı UC-011'dir).
- **Son Koşul:** Başarılı sipariş veritabanına kaydedilmiş, stok güncellenmiş ve müşteriye bildirim gönderilmiştir.
- **İlgili Gereksinim:** FR-3.5.1, FR-3.5.5, FR-3.3.2, FR-3.6.1, FR-2.1

---

### **Use Case 11: Ödeme Başarısız Olduğunda Hata Telafisi**

- **ID:** UC-011
- **Adı:** Ödeme Başarısız Olduğunda Hata Telafisi (Saga Deseni)
- **Aktör:** Müşteri
- **Ön Koşul:** Sipariş Servisi'nin ödeme isteği başarısız olmuştur (örn: Limit yetersizliği).
- **Temel Akış:**
    1. Müşteri, ödeme bilgilerini girer ve "Siparişi Tamamla" butonuna tıklar.
    2. Sipariş Servisi, ödeme isteğini Ödeme Sağlayıcısı'na gönderir.
    3. Ödeme Başarısız olur (örn: Hata Kodu 500).
    4. Sipariş Servisi, siparişin durumunu "Ödeme Başarısız" olarak günceller.
    5. Sipariş Servisi (Saga Deseni), eğer varsa önceden rezerve edilmiş kaynakları serbest bırakır (Bu akışta stok düşme olmadığı için telafi gerekmez).
    6. Ön Uç, müşteriyi "Ödeme Başarısız / Tekrar Dene" sayfasına yönlendirir.
- **Alternatif Akışlar:**
    - **6a. Stok Düşülüp Ödeme İptal Edilirse:** Ödeme alındıktan sonra stok düşülemediyse, Sipariş Servisi (Saga), Ödeme Servisi'ne İade/İptal isteği gönderir ve ödemeyi iade eder.
- **Son Koşul:** Müşteriye ödemenin başarısız olduğu bildirilmiş ve kaynaklar tutulmuyorsa, bir sonraki ödeme denemesi için hazırdır.
- **İlgili Gereksinim:** FR-3.5.7

---

### **Use Case 12: Kullanıcının Sipariş Durumunu Takip Etmesi**

- **ID:** UC-012
- **Adı:** Kullanıcının Sipariş Durumunu Takip Etmesi
- **Aktör:** Müşteri (Giriş Yapmış)
- **Ön Koşul:** Müşterinin en az bir tamamlanmış siparişi vardır.
- **Temel Akış:**
    1. Müşteri, "Siparişlerim" sayfasına gider.
    2. Sistem (API Gateway $\to$ Sipariş Servisi), kullanıcının tüm siparişlerinin listesini ister.
    3. Sipariş Servisi, sipariş geçmişini ve her siparişin güncel durumunu (örn: Kargoda) döndürür.
    4. Müşteri listeden bir siparişi seçer.
    5. Sipariş Servisi, detay sayfasında kargo takip numarasını ve lojistik durumunu (Kargoya Verildi, Teslim Edildi) gösterir.
- **Alternatif Akışlar:**
    - **3a. Sipariş Yok:** Müşterinin hiç siparişi yoksa, sistem "Henüz siparişiniz bulunmamaktadır" mesajı gösterir.
- **Son Koşul:** Müşteri, tüm siparişlerinin geçmişini ve anlık durumunu görebilir.
- **İlgili Gereksinim:** FR-3.5.6

---

### **Use Case 13: Admin Tarafından Ürün Oluşturma**

- **ID:** UC-013
- **Adı:** Admin Tarafından Yeni Ürün Oluşturma
- **Aktör:** ProductOwner (Admin Rolü)
- **Ön Koşul:** ProductOwner, Yönetim Paneline giriş yapmış ve gerekli yetkilere (RBAC) sahiptir.
- **Temel Akış:**
    1. ProductOwner, Yönetim Paneli üzerinden "Yeni Ürün Ekle" sayfasına gider.
    2. ProductOwner, ürün adı, açıklama, kategori, resimler ve varyant (renk, beden) bilgilerini girer.
    3. ProductOwner, başlangıç stok adedini de Stok Yönetim Servisi'ne girmek üzere belirtir.
    4. ProductOwner, "Kaydet" butonuna tıklar.
    5. Yönetim Paneli (API Gateway $\to$ Katalog Servisi'ne), Ürün oluşturma isteği gönderir.
    6. Katalog Servisi, veritabanına ürünü kaydeder ve Stok Servisi'ne (asenkron veya senkron) stok bilgisi gönderir.
- **Alternatif Akışlar:**
    - **3a. Yetkisizlik:** Eğer Admin yetkili değilse, sistem "Yetkiniz yok" hatası gösterir.
- **Son Koşul:** Yeni ürün, ürün kataloğuna eklenmiş ve stok bilgisi güncellenmiştir.
- **İlgili Gereksinim:** FR-3.2.1, FR-3.2.2, FR-3.7.1, FR-3.7.2, FR-3.1.6

---

### **Use Case 14: Admin Tarafından Sipariş Durumu Güncelleme**

- **ID:** UC-014
- **Adı:** Admin Tarafından Sipariş Durumu Güncelleme ve Bildirim
- **Aktör:** Operasyon Yöneticisi (Admin Rolü)
- **Ön Koşul:** Admin, yönetim panelinde bir siparişin detay sayfasındadır.
- **Temel Akış:**
    1. Admin, siparişin durumunu "Hazırlanıyor"dan "Kargoya Verildi" olarak değiştirmek için açılır menüyü kullanır.
    2. Admin, kargo takip numarasını girer.
    3. Admin, "Güncelle" butonuna tıklar.
    4. Yönetim Paneli, Sipariş Servisi'ne durum güncelleme isteği gönderir.
    5. Sipariş Servisi, sipariş durumunu günceller.
    6. Sipariş Servisi, "Sipariş Kargoya Verildi" olayını mesaj kuyruğuna yayınlar.
    7. Bildirim Servisi, olayı dinler ve müşteriye kargo takip numarası ile birlikte e-posta gönderir.
- **Alternatif Akışlar:**
    - **2a. Takip Numarası Eksik:** Admin takip numarasını girmezse, sistem uyarı gösterir ve güncelleme reddedilir.
- **Son Koşul:** Sipariş durumu güncellenmiş ve müşteri bilgilendirilmiştir.
- **İlgili Gereksinim:** FR-3.5.6, FR-3.6.1, FR-3.7.3

---

### **Use Case 15: İptal/İade Sonrası Stok Güncelleme**

- **ID:** UC-015
- **Adı:** İptal/İade Sonrası Stok Güncelleme
- **Aktör:** Operasyon Yöneticisi (veya Otomatik Sistem)
- **Ön Koşul:** Bir sipariş için İptal veya İade süreci onaylanmış ve sistem üzerinden yönetilmektedir.
- **Temel Akış:**
    1. Yönetim Paneli (veya harici sistem), Sipariş Servisi'ne bir siparişin "İade Tamamlandı" durumuna geçtiğini bildirir.
    2. Sipariş Servisi, "Sipariş İade Edildi" olayını mesajlaşma kuyruğuna yayınlar.
    3. Stok Yönetim Servisi, olayı dinler.
    4. Stok Servisi, iade edilen ürünlerin adedini mevcut stoka otomatik olarak ekler/günceller.
- **Alternatif Akışlar:**
    - **4a. Hatalı İade:** İade süreci bir hatadan dolayı yarıda kalırsa, telafi mekanizması devreye girer.
- **Son Koşul:** İade/İptal edilen ürünlerin stok adedi gerçek zamanlı olarak güncellenmiştir.
- **İlgili Gereksinim:** FR-3.3.3, FR-3.5.8, FR-2.1