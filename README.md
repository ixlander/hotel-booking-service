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
│   ├── app/               # Application entry point
│   │   └── main.go        # Main application logic
│   └── migrate/           # Migration utility
│       └── main.go        # Migration logic
│
├── docs/
│   ├── erd_diagram_hotel.png  # ERD diagram for hotel booking system
│   └── umls/                # UML sequence diagrams
│       ├── booking_cancel.puml
│       ├── check_available_hotels_and_rooms.puml
│       ├── create_booking.puml
│       └── user_login.puml
│
├── internal/
│   ├── app/
│   │   ├── config/        # Configuration (.env, etc.)
│   │   │   └── config.go  # Configuration loading
│   │   ├── connections/   # External connections (DB, APIs)
│   │   │   └── database.go # Database connection logic
│   │   ├── start/         # Startup helpers (routes, etc.)
│   │   │   └── routes.go  # Route initialization
│   │   └── store/         # Repository store
│   │       └── store.go   # Store handling
│   ├── data/              # Data layer (models)
│   │   ├── booking.go     # Booking data model
│   │   └── models.go      # General data models (hotel, user, etc.)
│   ├── deliveries/        # Delivery layer (HTTP handlers)
│   │   ├── http/          # HTTP delivery methods (controllers)
│   │   │   ├── auth_controller.go  # Auth-related handlers
│   │   │   ├── booking_controller.go  # Booking-related handlers
│   │   │   ├── hotel_controller.go  # Hotel-related handlers
│   │   │   └── user_controller.go   # User-related handlers
│   │   └── middleware/    # HTTP middleware (auth, logger)
│   │       ├── auth_middleware.go  # Auth validation middleware
│   │       └── logger.go  # Request logging middleware
│   ├── pkg/               # Internal libraries (utilities)
│   │   ├── apperror/      # Custom error handling
│   │   │   └── errors.go  # Error types and handling
│   │   ├── httputil/      # HTTP utilities
│   │   │   └── response.go # HTTP response formatting
│   │   ├── jwtutil/       # JWT utilities (signing and parsing)
│   │   │   └── jwt.go     # JWT token utility functions
│   │   └── logger/        # Logging utilities
│   │       └── logger.go  # Logger setup
│   ├── repositories/      # Repository layer (data access)
│   │   ├── booking_repository.go  # Booking repository
│   │   ├── hotel_repository.go    # Hotel repository
│   │   ├── room_repository.go     # Room repository
│   │   ├── user_repository.go     # User repository
│   │   └── postgres/        # Postgres repository
│   │       └── postgres.go  # Postgres connection and queries
│   ├── services/          # Service controllers
│   │   ├── booking_service.go  # Booking service logic
│   │   ├── hotel_service.go    # Hotel service logic
│   │   ├── room_service.go     # Room service logic
│   │   └── user_service.go     # User service logic
│   └── usecases/          # Business logic
│       ├── auth_usecase.go     # Auth use case logic
│       ├── booking_usecase.go  # Booking use case logic
│       ├── hotel_usecase.go    # Hotel use case logic
│       └── user_usecase.go     # User use case logic
│
├── migrations/            # Database migrations
│   ├── 000001_create_initial_tables.up.sql   # Initial table creation
│   ├── 000001_create_initial_tables.down.sql # Rollback for initial tables
│   ├── 000002_seed_data.up.sql    # Data seeding (up)
│   └── 000002_seed_data.down.sql  # Data seeding (down)
│
├── .env                   # Environment variables
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker build file
├── go.mod                 # Go module file
└── go.sum                 # Go dependencies checksum
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
