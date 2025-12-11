package usecase

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "simple-task-manager/internal/core/domain"
    "simple-task-manager/pkg/util"
)

type UserUsecase struct {
    userRepo       domain.UserRepository
    contextTimeout time.Duration
    jwtSecret      string // Secret key disuntikkan saat inisialisasi
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration, secret string) *UserUsecase {
    return &UserUsecase{
        userRepo:       userRepo,
        contextTimeout: timeout,
        jwtSecret:      secret,
    }
}

// Register user baru
func (u *UserUsecase) Register(c context.Context, user *domain.User) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    // 1. SANITISASI: Hapus spasi di awal/akhir nama, email, dan password
    user.Name = strings.TrimSpace(user.Name)
    user.Email = strings.TrimSpace(user.Email)
    cleanPassword := strings.TrimSpace(user.Password)

    // 2. Cek apakah email sudah ada
    existUser, err := u.userRepo.GetByEmail(ctx, user.Email)
    if err == nil && existUser != nil {
        return errors.New("email already exists")
    }

    // 3. Hash Password yang sudah dibersihkan
    hashedPass, err := util.HashPassword(cleanPassword)
    if err != nil {
        return err
    }
    user.Password = hashedPass

    // 4. Set timestamp
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    // 5. Simpan ke DB
    return u.userRepo.Create(ctx, user)
}

type LoginResponse struct {
    AccessToken string      `json:"access_token"`
    User        domain.User `json:"user"`
}

// Login user
func (u *UserUsecase) Login(c context.Context, email string, password string) (*LoginResponse, error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    // 1. SANITISASI INPUT LOGIN
    cleanEmail := strings.TrimSpace(email)
    cleanPassword := strings.TrimSpace(password)

    // 2. Cari user by email
    user, err := u.userRepo.GetByEmail(ctx, cleanEmail)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid email or password")
    }

    // --- DEBUGGING AREA ---
    fmt.Println("------------------------------------------------")
    fmt.Println("DEBUG LOGIN:")
    fmt.Printf("1. Email Input: %s\n", cleanEmail)
    fmt.Printf("2. Email di DB: %s\n", user.Email)
    fmt.Printf("3. Password Input: %s\n", cleanPassword)
    fmt.Printf("4. Password Hash di DB: %s\n", user.Password)
    fmt.Println("------------------------------------------------")
    // ----------------------

    // 3. Cek Password (gunakan cleanPassword)
    err = util.CheckPassword(cleanPassword, user.Password)
    if err != nil {
        fmt.Printf(">>> BCRYPT ERROR: %v\n", err)
        return nil, errors.New("invalid email or password")
    }

    // 4. Generate JWT Token (valid 24 jam)
    accessToken, err := util.CreateAccessToken(user.ID, u.jwtSecret, 24*time.Hour)
    if err != nil {
        return nil, err
    }

    // 5. Return response
    return &LoginResponse{
        AccessToken: accessToken,
        User:        *user,
    }, nil
}
