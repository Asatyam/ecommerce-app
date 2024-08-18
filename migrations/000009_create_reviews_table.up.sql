CREATE TABLE IF NOT EXISTS reviews(
    id bigint PRIMARY KEY,
    review text,
    rating int NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    customer_id bigint NOT NULL REFERENCES users,
    product_id bigint NOT NULL REFERENCES products ON DELETE CASCADE,
    version int NOT NULL DEFAULT 1
)