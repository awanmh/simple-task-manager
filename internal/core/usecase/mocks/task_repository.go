package mocks

import (
    "context"
    "simple-task-manager/internal/core/domain"

    "github.com/stretchr/testify/mock"
)

// TaskRepository adalah mock untuk domain.TaskRepository
type TaskRepository struct {
    mock.Mock
}

func (m *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *TaskRepository) Fetch(ctx context.Context, userID int64) ([]domain.Task, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *TaskRepository) Delete(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *TaskRepository) CreateSubtask(ctx context.Context, sub *domain.Subtask) error {
    args := m.Called(ctx, sub)
    return args.Error(0)
}

func (m *TaskRepository) DeleteSubtask(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *TaskRepository) ToggleSubtask(ctx context.Context, id int64) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}
