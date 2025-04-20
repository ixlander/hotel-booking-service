package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *data.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at`
	
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*data.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE id = $1`
	
	var user data.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*data.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`
	
	var user data.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}