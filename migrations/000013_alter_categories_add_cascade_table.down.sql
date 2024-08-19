ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS categories_parent_fkey;

ALTER TABLE categories
    ADD CONSTRAINT categories_parent_fkey
        FOREIGN KEY (parent) REFERENCES categories(id);