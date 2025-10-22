-- user/init.sql

-- Kullanıcı Rolleri için özel bir TYPE oluşturuluyor (PostgreSQL ENUM karşılığı)
CREATE TYPE user_role AS ENUM ('KULLANICI', 'SATICI', 'YONETICI');

-- Adres Tablosu
CREATE TABLE "Address"
(
    "AddressID"   BIGSERIAL PRIMARY KEY,
    "Country"     VARCHAR(255) NOT NULL,
    "City"        VARCHAR(255) NOT NULL,
    "District"    VARCHAR(255),
    "AddressLine" VARCHAR(500) NOT NULL,
    "PostCode"    VARCHAR(50),
    "CreateDate"  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "UpdateDate"  TIMESTAMP WITHOUT TIME ZONE,
    "DeleteDate"  TIMESTAMP WITHOUT TIME ZONE,
    "IsDeleted"   BOOLEAN      NOT NULL DEFAULT FALSE
);

-- Kullanıcı Tablosu
CREATE TABLE "User"
(
    "UserID"     BIGSERIAL PRIMARY KEY,
    "Email"      VARCHAR(255) UNIQUE NOT NULL,
    "Password"   VARCHAR(255)        NOT NULL,
    "Ad"         VARCHAR(255),
    "Soyad"      VARCHAR(255),
    "Role"       user_role           NOT NULL,
    "AddressID"  BIGINT,
    "CreateDate" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "UpdateDate" TIMESTAMP WITHOUT TIME ZONE,
    "DeleteDate" TIMESTAMP WITHOUT TIME ZONE,
    "IsDeleted"  BOOLEAN             NOT NULL DEFAULT FALSE,

    FOREIGN KEY ("AddressID") REFERENCES "Address" ("AddressID")
);

-- Adres Verileri
INSERT INTO "Address" ("Country", "City", "District", "AddressLine", "PostCode", "IsDeleted")
VALUES ('Türkiye', 'İstanbul', 'Kadıköy', 'Örnek Sokak, No: 5', '34710', FALSE),
       ('Türkiye', 'Ankara', 'Çankaya', 'Atatürk Bulvarı, 12/A', '06540', FALSE),
       ('ABD', 'New York', 'Manhattan', '5th Avenue, Apt 101', '10001', FALSE);

-- Kullanıcı Verileri (AddressID'leri yukarıdaki sırayla alacak şekilde)
INSERT INTO "User" ("Email", "Password", "Ad", "Soyad", "Role", "AddressID", "IsDeleted")
VALUES ('ali.demir@example.com', 'hashed_pass_1', 'Ali', 'Demir', 'KULLANICI', 1, FALSE),
       ('satıcı@magaza.com', 'hashed_pass_2', 'Elif', 'Yılmaz', 'SATICI', 2, FALSE),
       ('yonetici@sistem.com', 'hashed_pass_3', 'Ayşe', 'Kaya', 'YONETICI', 3, FALSE);