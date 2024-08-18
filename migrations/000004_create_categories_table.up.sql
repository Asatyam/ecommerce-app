CREATE TABLE IF NOT EXISTS categories(
    id bigint PRIMARY KEY,
    name text NOT NULL UNIQUE ,
    parent bigint NOT NULL DEFAULT -1 REFERENCES categories,
    version int NOT NULL DEFAULT 1
)