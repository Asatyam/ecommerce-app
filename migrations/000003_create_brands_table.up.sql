CREATE TABLE IF NOT EXISTS brands(
    id bigserial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    logo text NOT NULL,
    version int NOT NULL DEFAULT 1                           
)
