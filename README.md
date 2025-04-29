
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
│   └── migrate/           # Migration utility
│
├── internal/
│   ├── app/
│   │   ├── config/        # Configuration (.env, etc.)
│   │   ├── connections/   # External connections (DB, APIs)
│   │   ├── start/         # Startup helpers
│   │   └── store/         # Repository store
│   ├── data/              # Data layer (models)
│   ├── deliveries/        # Delivery layer (HTTP handlers)
│   ├── repositories/      # Repository layer (data access)
│   ├── services/          # Service controllers
│   ├── usecases/          # Business logic
│   └── pkg/               # Internal libraries
│
├── docker-compose.yml     # Docker Compose config
├── Dockerfile             # Docker build file
├── init.sql               # Database initialization
├── .env                   # Environment variables
├── go.mod                 # Go module file
└── go.sum                 # Go dependencies checksum
```

## Prerequisites

- Go 1.18 or higher
- PostgreSQL
- Docker & Docker Compose (optional, for containerized setup)

## Running Locally

### Option 1: Without Docker

1. Install Go:
   ```
   https://golang.org/doc/install
   ```

2. Install PostgreSQL and create a database:
   ```
   createdb hotel_booking
   ```

3. Clone the repository:
   ```
   git clone <repository-url>
   cd hotel-booking-service
   ```

4. Configure the environment variables:
   ```
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. Initialize the database:
   ```
   psql -U postgres -d hotel_booking -f init.sql
   ```

6. Install dependencies:
   ```
   go mod download
   ```

7. Build and run the application:
   ```
   go build -o hotel-booking-service ./cmd/app
   ./hotel-booking-service
   ```

8. The API will be available at `http://localhost:8080`

### Option 2: With Docker Compose

1. Install Docker and Docker Compose:
   ```
   https://docs.docker.com/get-docker/
   https://docs.docker.com/compose/install/
   ```

2. Clone the repository:
   ```
   git clone <repository-url>
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
  
  Request Body:
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- **Login**
  ```
  POST /login
  
  Request Body:
  {
    "email": "user@example.com",
    "password": "password123"
  }
  
  Response:
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
  
  Request Body:
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

1. Create a new SQL file in `cmd/migrate/migrations/` with a sequential number prefix
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