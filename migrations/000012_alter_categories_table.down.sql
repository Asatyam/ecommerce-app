ALTER TABLE categories
    ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE "categories_auto_increment";
