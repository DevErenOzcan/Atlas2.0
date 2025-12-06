-- Sipariş Tablosu (Order)
-- "Order" kelimesi PostgreSQL'de saklı kelime olduğu için çift tırnak ile kullanıldı.
CREATE TABLE "Order"
(
    "OrderID"    BIGSERIAL PRIMARY KEY,
    "UserID"     BIGINT         NOT NULL,
    "Total"      NUMERIC(10, 2) NOT NULL, -- FLOAT yerine NUMERIC kullanıldı
    "Final"      NUMERIC(10, 2),          -- Diyagramdaki BIGINT yerine NUMERIC varsayımı
    "IsShipped"  BOOLEAN        NOT NULL DEFAULT FALSE,
    "CreateDate" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "IsDeleted"  BOOLEAN        NOT NULL DEFAULT FALSE
    -- Ayrı veritabanı nedeniyle User tablosuna doğrudan FK tanımı yapılmıyor
);

-- Sipariş Detay Tablosu (OrderDetail)
CREATE TABLE "OrderDetail"
(
    "DetailID"   BIGSERIAL PRIMARY KEY,
    "OrderID"    BIGINT NOT NULL,
    "ProductID"  BIGINT NOT NULL,
    "Item"       INT    NOT NULL, -- Miktar (Quantity)
    "Final"      NUMERIC(10, 2),  -- Diyagramdaki BIGINT yerine NUMERIC varsayımı
    "CreateDate" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "IsDeleted"  BOOLEAN        NOT NULL DEFAULT FALSE
    -- Ayrı veritabanları nedeniyle Order ve Product tablolarına doğrudan FK tanımı yapılmıyor
);

INSERT INTO "Order" ("UserID", "Total", "Final", "IsShipped")
VALUES (1, 15000.00, 15000.00, TRUE), -- Ali Demir'in siparişi
       (1, 31.98, 31.98, FALSE),
       (3, 9.90, 9.90, TRUE);

-- Sipariş Detay Verileri
INSERT INTO "OrderDetail" ("OrderID", "ProductID", "Item", "Final")
VALUES (1, 1, 1, 15000.00),
       (2, 2, 2, 31.98),
       (3, 3, 1, 9.90);
