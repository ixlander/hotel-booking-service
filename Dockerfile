FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hotel-booking-service ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/hotel-booking-service .
COPY --from=builder /go/bin/migrate /usr/bin/migrate

COPY .env .env
COPY --from=builder /app/migrations /app/migrations

ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=hotel_booking \
    DB_SSLMODE=disable \
    SERVER_PORT=8080 \
    JWT_SECRET=your_secret_key_replace_this_in_production \
    JWT_TOKEN_EXPIRY_HOURS=24

EXPOSE 8080

CMD ["/app/hotel-booking-service"]