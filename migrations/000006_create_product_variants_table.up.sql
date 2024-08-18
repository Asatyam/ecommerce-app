CREATE TABLE IF NOT EXISTS product_variants(
    id bigserial PRIMARY KEY,
    price int NOT NULL,
    discount int NOT NULL DEFAULT 0,
    sku text NOT NULL UNIQUE,
    variants jsonb NOT NULL,
    product_id bigint NOT NULL REFERENCES products ON DELETE CASCADE,
    version int NOT NULL DEFAULT 1
)