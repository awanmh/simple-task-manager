package mocks

import (
    "context"
    "simple-task-manager/internal/core/domain"

    "github.com/stretchr/testify/mock"
)

// UserRepository adalah mock untuk domain.UserRepository
type UserRepository struct {
    mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    args := m.Called(ctx, email)
    if u := args.Get(0); u != nil {
        return u.(*domain.User), args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    args := m.Called(ctx, id)
    if u := args.Get(0); u != nil {
        return u.(*domain.User), args.Error(1)
    }
    return nil, args.Error(1)
}
