-- Kategori Tablosu
CREATE TABLE "Category"
(
    "CatID"   BIGSERIAL PRIMARY KEY,
    "CatName" VARCHAR(255) UNIQUE NOT NULL
);

-- Para Birimi Tablosu
CREATE TABLE "Currency"
(
    "Alt"     VARCHAR(255) PRIMARY KEY,
    "AltName" VARCHAR(255) NOT NULL
);

-- Ürün Tablosu
CREATE TABLE "Product"
(
    "ProductID"     BIGSERIAL PRIMARY KEY,
    "Alt"           VARCHAR(255)   NOT NULL, -- Para Birimi FK
    "Category"      BIGINT         NOT NULL, -- Kategori FK
    "ProductName"   VARCHAR(255)   NOT NULL,
    "Description"   TEXT,
    "DimensDetails" TEXT,                    -- "Dimens Details" sütunu birleştirildi
    "Stock"         INT            NOT NULL,
    "Price"         NUMERIC(10, 2) NOT NULL, -- FLOAT yerine daha hassas NUMERIC kullanıldı
    "Currency"      VARCHAR(255)   NOT NULL, -- Para birimi alanının tekrarı
    "SellerID"      BIGINT,                  -- Satıcı (user servisi ID'si)
    "CreateDate"    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "UpdateDate"    TIMESTAMP WITHOUT TIME ZONE,
    "DeleteDate"    TIMESTAMP WITHOUT TIME ZONE,
    "IsDeleted"     BOOLEAN        NOT NULL DEFAULT FALSE,

    FOREIGN KEY ("Alt") REFERENCES "Currency" ("Alt"),
    FOREIGN KEY ("Category") REFERENCES "Category" ("CatID")
    -- Her servis kendi veritabanına sahip olduğundan, User tablosuna doğrudan FK tanımlanmaz
);


INSERT INTO "Currency" ("Alt", "AltName") VALUES
                                              ('TRY', 'Türk Lirası'),
                                              ('USD', 'Amerikan Doları'),
                                              ('EUR', 'Euro');

-- Kategori Verileri
INSERT INTO "Category" ("CatName")
VALUES ('Elektronik'),
       ('Giyim'),
       ('Kitap');

-- Ürün Verileri (Category: 1=Elektronik, 2=Giyim, 3=Kitap. SellerID: 2=Elif Yılmaz)
INSERT INTO "Product" ("Alt", "Category", "ProductName", "Description", "DimensDetails", "Stock", "Price", "Currency",
                       "SellerID", "IsDeleted")
VALUES ('TRY', 1, 'Akıllı Telefon X', 'Yüksek performanslı akıllı telefon.', '15x7x0.8 cm', 50, 15000.00, 'TRY', 2,
        FALSE),
       ('USD', 2, 'Pamuklu T-shirt', 'Yazlık, %100 pamuk.', 'S, M, L Bedenler', 120, 15.99, 'USD', 2, FALSE),
       ('EUR', 3, 'Bilim Kurgu Romanı', 'Çok satan bir bilim kurgu eseri.', '20x13x3 cm, Ciltli', 30, 9.90, 'EUR', 2,
        FALSE);