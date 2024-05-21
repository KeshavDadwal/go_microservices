CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    transaction_date TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
