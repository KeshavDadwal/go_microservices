CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) UNIQUE NOT NULL
);

INSERT INTO users (name, email, password) VALUES
('Alice Johnson', 'alice@example.com', 'hashedpassword1'),
('Bob Smith', 'bob@example.com', 'hashedpassword2'),
('John Smith', 'john@example.com', 'hashedpassword3'),
('Peter Smith', 'peter@example.com', 'hashedpassword4'),
('Mark Smith', 'mark@example.com', 'hashedpassword5'),
('Steve Smith', 'steve@example.com', 'hashedpassword6');
