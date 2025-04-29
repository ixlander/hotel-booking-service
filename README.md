# Hotel Booking Service

A RESTful API service for hotel room booking built with Go, following clean architecture principles.

## Features

- User registration and authentication with JWT
- Hotel and room listing with availability filtering
- Room booking creation and cancellation
- User booking history

## Project Structure

```
/hotel-booking-service
│
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── migrate/
│       └── main.go
│
├── docs/
│   ├── erd_diagram_hotel.png
│   └── umls/
│       ├── booking_cancel.puml
│       ├── check_available_hotels_and_rooms.puml
│       ├── create_booking.puml
│       └── user_login.puml
│
├── internal/
│   ├── app/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── connections/
│   │   │   └── database.go
│   │   ├── start/
│   │   │   └── routes.go
│   │   └── store/
│   │       └── store.go
│   ├── data/
│   │   ├── booking.go
│   │   └── models.go
│   ├── deliveries/
│   │   ├── http/
│   │   │   ├── auth_controller.go
│   │   │   ├── booking_controller.go
│   │   │   ├── hotel_controller.go
│   │   │   └── user_controller.go
│   │   └── middleware/
│   │       ├── auth_middleware.go
│   │       └── logger.go
│   ├── pkg/
│   │   ├── apperror/
│   │   │   └── errors.go
│   │   ├── httputil/
│   │   │   └── response.go
│   │   ├── jwtutil/
│   │   │   └── jwt.go
│   │   └── logger/
│   │       └── logger.go
│   ├── repositories/
│   │   ├── booking_repository.go
│   │   ├── hotel_repository.go
│   │   ├── room_repository.go
│   │   ├── user_repository.go
│   │   └── postgres/
│   │       └── postgres.go
│   ├── services/
│   │   ├── booking_service.go
│   │   ├── hotel_service.go
│   │   ├── room_service.go
│   │   └── user_service.go
│   └── usecases/
│       ├── auth_usecase.go
│       ├── booking_usecase.go
│       ├── hotel_usecase.go
│       └── user_usecase.go
│
├── migrations/
│   ├── 000001_create_initial_tables.up.sql
│   ├── 000001_create_initial_tables.down.sql
│   ├── 000002_seed_data.up.sql
│   └── 000002_seed_data.down.sql
│
├── .env
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Running Locally With Docker Compose

1. Install Docker and Docker Compose:
   ```
   https://docs.docker.com/get-docker/
   https://docs.docker.com/compose/install/
   ```

2. Clone the repository:
   ```
   git clone https://github.com/ixlander/hotel-booking-service.git
   cd hotel-booking-service
   ```

3. Start the application:
   ```
   docker-compose up -d
   ```

4. The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- **Register a new user**
  ```
  POST /register
  ```

  Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- **Login**
  ```
  POST /login
  ```

  Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

  Response:
  ```json
  {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "created_at": "2023-01-01T00:00:00Z"
    }
  }
  ```

### Hotels

- **Get all hotels with available rooms**
  ```
  GET /hotels?from_date=2023-01-01&to_date=2023-01-05
  ```

- **Get a specific hotel with available rooms**
  ```
  GET /hotels/1?from_date=2023-01-01&to_date=2023-01-05
  ```

### Bookings (Protected Routes - Require Authentication)

For these endpoints, include the JWT token in the Authorization header:
```
Authorization: Bearer your-jwt-token
```

- **Create a booking**
  ```
  POST /bookings
  ```

  Request Body:
  ```json
  {
    "room_id": 1,
    "from_date": "2023-01-01T00:00:00Z",
    "to_date": "2023-01-05T00:00:00Z"
  }
  ```

- **Get user bookings**
  ```
  GET /bookings
  ```

- **Cancel a booking**
  ```
  DELETE /bookings/1
  ```

## Development

### Adding Database Migrations

The project is set up to use SQL migrations. To add a new migration:

1. Create a new SQL file in `migrations/` with a sequential number prefix
2. Write your migration SQL (both up and down)
3. Run the migration tool:
   ```
   go run ./cmd/migrate
   ```

### Testing

Run the tests with:
```
go test ./...
```

## Deployment

For production deployment, make sure to:

1. Use a secure JWT secret key
2. Configure proper database credentials
3. Use SSL for database connection
4. Consider using a reverse proxy like Nginx

## License

This project is licensed under the MIT License - see the LICENSE file for details.
