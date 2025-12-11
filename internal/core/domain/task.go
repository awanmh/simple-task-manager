package domain

import (
    "context"
    "time"
)

// Subtask merepresentasikan checklist item di dalam sebuah Task
type Subtask struct {
    ID        int64     `json:"id"`
    TaskID    int64     `json:"task_id"` // Foreign Key ke Task
    Title     string    `json:"title"`
    IsDone    bool      `json:"is_done"`
    CreatedAt time.Time `json:"created_at"`
}

// Task merepresentasikan tugas utama
type Task struct {
    ID           int64      `json:"id"`
    UserID       int64      `json:"user_id"`       // Foreign Key ke User
    Title        string     `json:"title"`
    Description  string     `json:"description"`
    Status       string     `json:"status"`        // e.g., "pending", "in_progress", "done"
    Priority     string     `json:"priority"`      // low, medium, high
    Labels       []string   `json:"labels"`        // contoh: ["work", "bug"]
    ReminderTime *time.Time `json:"reminder_time"` // pointer agar bisa null

    // --- FIELD BARU ---
    RecurrencePattern string     `json:"recurrence_pattern"` // "daily", "weekly", "monthly"
    NextRun           *time.Time `json:"next_run"`           // kapan trigger berikutnya
    Subtasks          []Subtask  `json:"subtasks"`           // Nested JSON untuk checklist
    // ------------------

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// TaskRepository mendefinisikan kontrak untuk operasi database terkait Task & Subtask
type TaskRepository interface {
    // --- Method Task ---
    Create(ctx context.Context, task *Task) error
    Fetch(ctx context.Context, userID int64) ([]Task, error)
    GetByID(ctx context.Context, id int64) (*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id int64) error

    // --- Method Subtask ---
    CreateSubtask(ctx context.Context, sub *Subtask) error
    DeleteSubtask(ctx context.Context, id int64) error
    ToggleSubtask(ctx context.Context, id int64) error
}
