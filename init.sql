CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hotels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL
);

CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    hotel_id INT REFERENCES hotels(id) ON DELETE CASCADE,
    number VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'confirmed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO hotels (name, city) VALUES
('Grand Hotel', 'New York'),
('Seaside Resort', 'Miami'),
('Mountain Lodge', 'Denver');

INSERT INTO rooms (hotel_id, number, capacity, price) VALUES
(1, '101', 2, 150.00),
(1, '102', 2, 150.00),
(1, '103', 4, 250.00),
(2, '201', 2, 200.00),
(2, '202', 3, 250.00),
(3, '301', 2, 180.00),
(3, '302', 4, 280.00);