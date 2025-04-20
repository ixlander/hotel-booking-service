INSERT INTO users (email, password) VALUES 
('user1@example.com', '$2a$10$1JD4XYBZRJd5IEjW4ReECOe77GgD57NUu8JyQmYr5qWOlYMYvkwma'),
('user2@example.com', '$2a$10$1JD4XYBZRJd5IEjW4ReECOe77GgD57NUu8JyQmYr5qWOlYMYvkwma');

INSERT INTO hotels (name, city) VALUES
('Grand Hotel', 'New York'),
('Sea View Resort', 'Miami'),
('Mountain Retreat', 'Denver');

INSERT INTO rooms (hotel_id, number, capacity, price) VALUES
(1, '101', 2, 150.00),
(1, '102', 2, 150.00),
(1, '201', 4, 250.00),
(1, '301', 1, 100.00);

INSERT INTO rooms (hotel_id, number, capacity, price) VALUES
(2, 'A1', 2, 200.00),
(2, 'A2', 2, 200.00),
(2, 'B1', 3, 300.00),
(2, 'P1', 4, 500.00);

INSERT INTO rooms (hotel_id, number, capacity, price) VALUES
(3, 'R1', 2, 175.00),
(3, 'R2', 2, 175.00),
(3, 'S1', 4, 350.00);

INSERT INTO bookings (user_id, room_id, from_date, to_date, status) VALUES
(1, 1, '2025-05-01', '2025-05-05', 'confirmed'),
(2, 5, '2025-06-10', '2025-06-15', 'confirmed');
