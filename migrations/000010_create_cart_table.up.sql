CREATE TABLE IF NOT EXISTS cart(
    id bigserial PRIMARY KEY,
    customer_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    variant_id bigint NOT NULL REFERENCES product_variants ON DELETE CASCADE,
    quantity int NOT NULL,
    version bigint NOT NULL  DEFAULT 1
)