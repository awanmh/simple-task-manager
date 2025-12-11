package mocks

import (
    "context"
    "simple-task-manager/internal/core/domain"
)

type UserRepositoryMock struct {
    CreateFn     func(ctx context.Context, user *domain.User) error
    GetByEmailFn func(ctx context.Context, email string) (*domain.User, error)
    GetByIDFn    func(ctx context.Context, id int64) (*domain.User, error)
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
    if m.CreateFn != nil {
        return m.CreateFn(ctx, user)
    }
    return nil
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    if m.GetByEmailFn != nil {
        return m.GetByEmailFn(ctx, email)
    }
    return nil, nil
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    if m.GetByIDFn != nil {
        return m.GetByIDFn(ctx, id)
    }
    return nil, nil
}
