CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    prime bool NOT NULL DEFAULT false,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    is_admin bool NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);