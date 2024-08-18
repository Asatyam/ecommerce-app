CREATE TABLE IF NOT EXISTS addresses(
    id bigserial PRIMARY KEY,
    customer_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    address_line text NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    pin text NOT NULL,
    version int NOT NULL DEFAULT 1
)