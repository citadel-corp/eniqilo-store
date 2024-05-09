CREATE TYPE product_categories AS ENUM('Clothing', 'Accessories', 'Footwear', 'Beverages');

CREATE TABLE IF NOT EXISTS
products (
    id VARCHAR(16) PRIMARY KEY,
    user_id VARCHAR(16) NOT NULL,
    name VARCHAR(30) NOT NULL,
    sku VARCHAR(30) NOT NULL,
    category product_categories NOT NULL,
    image_url TEXT NOT NULL,
    notes VARCHAR(200) NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL,
    location VARCHAR(200) NOT NULL,
    is_available boolean NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE products
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);

CREATE INDEX IF NOT EXISTS products_name
	ON products USING HASH(lower(name));
CREATE INDEX IF NOT EXISTS products_user_id
	ON products (user_id);
CREATE INDEX IF NOT EXISTS products_sku
	ON products USING HASH(sku);
CREATE INDEX IF NOT EXISTS products_category
	ON products (category);
CREATE INDEX IF NOT EXISTS products_is_available
	ON products (is_available);
CREATE INDEX IF NOT EXISTS products_in_stock
	ON products (stock) WHERE stock > 0;
CREATE INDEX IF NOT EXISTS products_not_in_stock
	ON products (stock) WHERE stock = 0;
CREATE INDEX IF NOT EXISTS products_created_at_desc
	ON products(created_at DESC);
CREATE INDEX IF NOT EXISTS products_created_at_asc
	ON products(created_at ASC);
CREATE INDEX IF NOT EXISTS products_price_desc
	ON products(price DESC);
CREATE INDEX IF NOT EXISTS products_price_asc
	ON products(price ASC);