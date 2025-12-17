package domain

import (
	"context"
	"time"
)
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	GetByID(ctx context.Context, id int64) (*User, error)
}