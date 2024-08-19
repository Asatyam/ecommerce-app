CREATE SEQUENCE categories_auto_increment;

ALTER TABLE categories
    ALTER COLUMN id SET DEFAULT nextval('categories_auto_increment');


