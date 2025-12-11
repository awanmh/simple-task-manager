package util

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// CreateAccessToken membuat JWT token baru
func CreateAccessToken(userID int64, secretKey string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// Catatan: Fungsi validasi token akan kita pakai nanti di Middleware.