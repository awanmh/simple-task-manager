package repository

import (
	"context"
	"errors"

	"simple-task-manager/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository adalah factory function untuk inisialisasi repository
func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Kita menggunakan QueryRow karena kita mengharapkan return ID dari database
	err := r.db.QueryRow(ctx, query, 
		user.Name, 
		user.Email, 
		user.Password, 
		user.CreatedAt, 
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User tidak ditemukan, return nil tanpa error sistem
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}