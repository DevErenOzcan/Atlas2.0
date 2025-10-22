-- MySQL exporter kullanıcısını varsa sil ve yeniden oluştur
DROP USER IF EXISTS '{{ .Values.secrets.exporterUser }}'@'%';
CREATE USER '{{ .Values.secrets.exporterUser }}'@'%' IDENTIFIED BY '{{ .Values.secrets.exporterPassword }}';
FLUSH PRIVILEGES;

USE {{ .Values.config.dbName }};

-- users tablosunu oluştur
CREATE TABLE users
(
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username     VARCHAR(255) NOT NULL UNIQUE,
    password     TEXT         NOT NULL,
    is_superuser BOOLEAN   DEFAULT FALSE,
    is_deleted   BOOLEAN   DEFAULT FALSE,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_date TIMESTAMP NULL DEFAULT NULL,
    INDEX        idx_users_username (username)
);

-- todos tablosunu oluştur
CREATE TABLE todos
(
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name          VARCHAR(500) NOT NULL,
    user_id       BIGINT UNSIGNED NOT NULL,
    complete_perc DECIMAL(5, 2) DEFAULT 0.00,
    is_deleted    BOOLEAN       DEFAULT FALSE,
    created_date  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    updated_date  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_date  TIMESTAMP NULL DEFAULT NULL,
    INDEX         idx_todos_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- todo_steps tablosunu oluştur
CREATE TABLE todo_steps
(
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    todo_id      BIGINT UNSIGNED NOT NULL,
    name         VARCHAR(500) NOT NULL,
    is_completed BOOLEAN   DEFAULT FALSE,
    is_deleted   BOOLEAN   DEFAULT FALSE,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_date TIMESTAMP NULL DEFAULT NULL,
    INDEX        idx_todo_steps_todo_id (todo_id),
    FOREIGN KEY (todo_id) REFERENCES todos (id) ON DELETE CASCADE
);

-- Kullanıcıları ekle
INSERT INTO users (username, password, is_superuser)
VALUES ('erenozcan', '123456', FALSE),
       ('privia', '123456', FALSE),
       ('admin', '123456', TRUE);

-- erenozcan (user_id: 1) için todoları ekle
INSERT INTO todos (name, user_id, complete_perc)
VALUES ('mülakata gir', 1, 0.0),
       ('işe git', 1, 40.0),
       ('kod yaz', 1, 100.0),
       ('proje sunumu', 1, 0.0);

-- privia (user_id: 2) için todoları ekle
INSERT INTO todos (name, user_id, complete_perc)
VALUES ('toplantı', 2, 0.0),
       ('rapor yaz', 2, 20.0),
       ('veri analizi yap', 2, 100.0),
       ('sunum hazırla', 2, 60.0),
       ('yazılım test et', 2, 0.0);

-- "mülakata gir" (todo_id: 1) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (1, 'özgeçmiş güncelle', TRUE),
       (1, 'şirketi araştır', TRUE),
       (1, 'mülakat sorularını çalış', FALSE),
       (1, 'kıyafet hazırla', FALSE),
       (1, 'mülakata git', FALSE);

-- "işe git" (todo_id: 2) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (2, 'alarm kur', TRUE),
       (2, 'kahvaltı yap', TRUE),
       (2, 'işe gitmek için yola çık', FALSE),
       (2, 'günlük toplantıya katıl', FALSE),
       (2, 'görevleri tamamla', FALSE);

-- "kod yaz" (todo_id: 3) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (3, 'gereksinimleri analiz et', TRUE),
       (3, 'tasarım yap', TRUE),
       (3, 'kodu yaz', TRUE),
       (3, 'test et', TRUE),
       (3, 'kod incelemesi yap', TRUE);

-- "proje sunumu" (todo_id: 4) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (4, 'sunum içeriğini hazırla', FALSE),
       (4, 'slaytları oluştur', FALSE),
       (4, 'prova yap', FALSE),
       (4, 'sunumu gerçekleştir', FALSE),
       (4, 'geri bildirimleri topla', FALSE);

-- "toplantı" (todo_id: 5) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (5, 'katılımcı listesi oluştur', FALSE),
       (5, 'sunum için hazırlık yap', FALSE),
       (5, 'toplantı hedeflerini belirle', FALSE),
       (5, 'katılımcılara toplantı notlarını gönder', FALSE),
       (5, 'toplantıdan sonra özet çıkar', FALSE);

-- "rapor yaz" (todo_id: 6) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (6, 'başlıkları belirle', TRUE),
       (6, 'verileri topla', FALSE),
       (6, 'görselleri hazırlayın', FALSE),
       (6, 'yazmaya başla', FALSE),
       (6, 'son kontrolleri yap', FALSE);

-- "veri analizi yap" (todo_id: 7) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (7, 'veri setini incele', TRUE),
       (7, 'veri temizliği yap', TRUE),
       (7, 'istatistiksel analiz yap', TRUE),
       (7, 'veri görselleştirme yap', TRUE),
       (7, 'sonuçları yorumla', TRUE);

-- "sunum hazırla" (todo_id: 8) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (8, 'konu başlıklarını belirle', TRUE),
       (8, 'görselleri seç', TRUE),
       (8, 'slaytları oluştur', TRUE),
       (8, 'sunum metni hazırla', FALSE),
       (8, 'son kontrol yap', FALSE);

-- "yazılım test et" (todo_id: 9) için adımları ekle
INSERT INTO todo_steps (todo_id, name, is_completed)
VALUES (9, 'test senaryolarını yaz', FALSE),
       (9, 'yazılımı kur', FALSE),
       (9, 'fonksiyonları test et', FALSE),
       (9, 'hata ayıklama yap', FALSE),
       (9, 'sonuçları raporla', FALSE);