CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    bio TEXT
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    genres TEXT, 
    published_at TIMESTAMP NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL,

    author_id INT NOT NULL REFERENCES authors(id) ON DELETE CASCADE
);



CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,

    UNIQUE (street, city, state, postal_code, country)
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    address_id INT NOT NULL REFERENCES addresses(id) ON DELETE CASCADE
);


CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    total_price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) NOT NULL
);


CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    quantity INT NOT NULL
);
