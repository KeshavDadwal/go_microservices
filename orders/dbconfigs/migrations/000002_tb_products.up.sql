CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL 
);

INSERT INTO products (name, description, price) VALUES
('Laptop', 'High performance laptop.', 999.99),
('Smartphone', 'Latest model smartphone.', 399.99),
('I pad', 'High performance laptop.', 9929.99),
('Iphone', 'Latest model smartphone.', 299.99),
('Asus Laptop', 'High performance laptop.', 9929.99),
('Adnroid Phone', 'Latest model smartphone.', 2399.99);
