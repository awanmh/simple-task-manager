package domain

import (
	"context"
	"time"
)

// User merepresentasikan data pengguna dalam sistem
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	// Password tidak di-export ke JSON demi keamanan
	Password  string    `json:"-"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository mendefinisikan kontrak untuk operasi database terkait User
type UserRepository interface {
	// Create menyimpan user baru ke database
	Create(ctx context.Context, user *User) error
	
	// GetByEmail mencari user berdasarkan email (berguna untuk login)
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// GetByID mencari user berdasarkan primary key
	GetByID(ctx context.Context, id int64) (*User, error)
}