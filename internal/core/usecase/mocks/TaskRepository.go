package mocks

import (
    "context"
    "simple-task-manager/internal/core/domain"
)

// TaskRepositoryMock adalah implementasi manual dari domain.TaskRepository
// dengan function field agar fleksibel di unit test.
type TaskRepositoryMock struct {
    CreateFn        func(ctx context.Context, task *domain.Task) error
    FetchFn         func(ctx context.Context, userID int64) ([]domain.Task, error)
    GetByIDFn       func(ctx context.Context, id int64) (*domain.Task, error)
    UpdateFn        func(ctx context.Context, task *domain.Task) error
    DeleteFn        func(ctx context.Context, id int64) error
    CreateSubtaskFn func(ctx context.Context, sub *domain.Subtask) error
    DeleteSubtaskFn func(ctx context.Context, id int64) error
    ToggleSubtaskFn func(ctx context.Context, id int64) error
}

func (m *TaskRepositoryMock) Create(ctx context.Context, task *domain.Task) error {
    if m.CreateFn != nil {
        return m.CreateFn(ctx, task)
    }
    return nil
}

func (m *TaskRepositoryMock) Fetch(ctx context.Context, userID int64) ([]domain.Task, error) {
    if m.FetchFn != nil {
        return m.FetchFn(ctx, userID)
    }
    return nil, nil
}

func (m *TaskRepositoryMock) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
    if m.GetByIDFn != nil {
        return m.GetByIDFn(ctx, id)
    }
    return nil, nil
}

func (m *TaskRepositoryMock) Update(ctx context.Context, task *domain.Task) error {
    if m.UpdateFn != nil {
        return m.UpdateFn(ctx, task)
    }
    return nil
}

func (m *TaskRepositoryMock) Delete(ctx context.Context, id int64) error {
    if m.DeleteFn != nil {
        return m.DeleteFn(ctx, id)
    }
    return nil
}

func (m *TaskRepositoryMock) CreateSubtask(ctx context.Context, sub *domain.Subtask) error {
    if m.CreateSubtaskFn != nil {
        return m.CreateSubtaskFn(ctx, sub)
    }
    return nil
}

func (m *TaskRepositoryMock) DeleteSubtask(ctx context.Context, id int64) error {
    if m.DeleteSubtaskFn != nil {
        return m.DeleteSubtaskFn(ctx, id)
    }
    return nil
}

func (m *TaskRepositoryMock) ToggleSubtask(ctx context.Context, id int64) error {
    if m.ToggleSubtaskFn != nil {
        return m.ToggleSubtaskFn(ctx, id)
    }
    return nil
}
