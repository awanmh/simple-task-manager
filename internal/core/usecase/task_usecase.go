package usecase

import (
    "context"
    "errors"
    "time"

    "simple-task-manager/internal/core/domain"
)

type TaskUsecase struct {
    taskRepo       domain.TaskRepository
    contextTimeout time.Duration
}

func NewTaskUsecase(taskRepo domain.TaskRepository, timeout time.Duration) *TaskUsecase {
    return &TaskUsecase{
        taskRepo:       taskRepo,
        contextTimeout: timeout,
    }
}

// --- TASK METHODS ---

// Create task baru
func (u *TaskUsecase) Create(c context.Context, task *domain.Task) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()

    if task.Status == "" {
        task.Status = "pending"
    }

    return u.taskRepo.Create(ctx, task)
}

// Fetch semua task milik user tertentu
func (u *TaskUsecase) Fetch(c context.Context, userID int64) ([]domain.Task, error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    return u.taskRepo.Fetch(ctx, userID)
}

// Update status task
func (u *TaskUsecase) UpdateStatus(c context.Context, id int64, userID int64, status string) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    task, err := u.taskRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if task == nil {
        return errors.New("task not found")
    }

    if task.UserID != userID {
        return errors.New("unauthorized: you don't own this task")
    }

    task.Status = status
    task.UpdatedAt = time.Now()

    return u.taskRepo.Update(ctx, task)
}

// Delete task
func (u *TaskUsecase) Delete(c context.Context, id int64, userID int64) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    task, err := u.taskRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if task == nil {
        return errors.New("task not found")
    }

    if task.UserID != userID {
        return errors.New("unauthorized: you don't own this task")
    }

    return u.taskRepo.Delete(ctx, id)
}

// --- SUBTASK METHODS ---

func (u *TaskUsecase) AddSubtask(c context.Context, taskID int64, title string) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    sub := &domain.Subtask{
        TaskID: taskID,
        Title:  title,
        IsDone: false,
    }
    return u.taskRepo.CreateSubtask(ctx, sub)
}

func (u *TaskUsecase) ToggleSubtask(c context.Context, id int64) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    return u.taskRepo.ToggleSubtask(ctx, id)
}

func (u *TaskUsecase) DeleteSubtask(c context.Context, id int64) error {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    return u.taskRepo.DeleteSubtask(ctx, id)
}
