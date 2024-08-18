CREATE TABLE IF NOT EXISTS products(
    id bigserial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    brand_id bigint NOT NULL REFERENCES brands ON DELETE CASCADE,
    category_id bigint NOT NULL REFERENCES categories ON DELETE CASCADE,
    version int NOT NULL DEFAULT 1
)