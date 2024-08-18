CREATE TABLE IF NOT EXISTS orders(
    id bigserial PRIMARY KEY,
    status text NOT NULL DEFAULT 'Confirmed',
    payment_status text NOT NULL DEFAULT 'NOT PAID',
    total int NOT NULL,
    contact_no varchar(20) NOT NULL,
    date DATE DEFAULT CURRENT_DATE,
    customer_id bigint NOT NULL REFERENCES users,
    address text NOT NULL,
    version int NOT NULL DEFAULT 1
)