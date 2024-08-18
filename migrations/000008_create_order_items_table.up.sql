CREATE TABLE IF NOT EXISTS order_items(
    id bigserial PRIMARY KEY,
    quantity int NOT NULL,
    price int NOT NULL,
    order_id bigint NOT NULL REFERENCES orders ON DELETE CASCADE,
    variant_id bigint NOT NULL REFERENCES product_variants,
    version int NOT NULL DEFAULT 1
)