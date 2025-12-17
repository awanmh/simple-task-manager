package repository

import (
    "context"
    "simple-task-manager/internal/core/domain"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskRepository struct {
    db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) domain.TaskRepository {
    return &PostgresTaskRepository{db: db}
}

// Create task baru
func (r *PostgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
    query := `
        INSERT INTO tasks (
            user_id, title, description, status, priority, labels, reminder_time, recurrence_pattern, next_run, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `

    // Set default priority jika kosong
    if task.Priority == "" {
        task.Priority = "medium"
    }

    err := r.db.QueryRow(ctx, query,
        task.UserID,
        task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.Labels,
        task.ReminderTime,
        task.RecurrencePattern,
        task.NextRun,
        task.CreatedAt,
        task.UpdatedAt,
    ).Scan(&task.ID)

    return err
}

// --- SUBTASK METHODS ---

func (r *PostgresTaskRepository) CreateSubtask(ctx context.Context, sub *domain.Subtask) error {
    query := `INSERT INTO subtasks (task_id, title, is_done, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
    return r.db.QueryRow(ctx, query, sub.TaskID, sub.Title, sub.IsDone).Scan(&sub.ID)
}

func (r *PostgresTaskRepository) DeleteSubtask(ctx context.Context, id int64) error {
    _, err := r.db.Exec(ctx, "DELETE FROM subtasks WHERE id = $1", id)
    return err
}

func (r *PostgresTaskRepository) ToggleSubtask(ctx context.Context, id int64) error {
    _, err := r.db.Exec(ctx, "UPDATE subtasks SET is_done = NOT is_done WHERE id = $1", id)
    return err
}

// --- TASK METHODS ---

// Fetch semua task milik user (lengkap dengan subtasks)
func (r *PostgresTaskRepository) Fetch(ctx context.Context, userID int64) ([]domain.Task, error) {
    query := `
        SELECT 
            t.id, 
            t.user_id, 
            t.title, 
            COALESCE(t.description, ''), 
            t.status, 
            COALESCE(t.priority, 'medium'), 
            t.labels, 
            t.reminder_time,
            
            -- FITUR BARU: Handle NULL Recurrence
            COALESCE(t.recurrence_pattern, ''), 
            t.next_run,
            
            t.created_at, 
            t.updated_at,
            COALESCE(
                json_agg(
                    json_build_object('id', s.id, 'task_id', s.task_id, 'title', s.title, 'is_done', s.is_done)
                ) FILTER (WHERE s.id IS NOT NULL), 
                '[]'
            ) as subtasks
        FROM tasks t
        LEFT JOIN subtasks s ON s.task_id = t.id
        WHERE t.user_id = $1
        GROUP BY t.id
        ORDER BY t.created_at DESC
    `

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    tasks := []domain.Task{}
    for rows.Next() {
        var t domain.Task
        err := rows.Scan(
            &t.ID, &t.UserID, &t.Title, &t.Description, &t.Status,
            &t.Priority, &t.Labels, &t.ReminderTime,
            &t.RecurrencePattern,
            &t.NextRun,
            &t.CreatedAt, &t.UpdatedAt,
            &t.Subtasks,
        )
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }

    return tasks, nil
}

// GetByID task berdasarkan ID
func (r *PostgresTaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
    query := `
        SELECT 
            t.id, t.user_id, t.title, 
            COALESCE(t.description, ''), 
            t.status, 
            COALESCE(t.priority, 'medium'), 
            t.labels, t.reminder_time,
            
            -- FITUR BARU
            COALESCE(t.recurrence_pattern, ''), 
            t.next_run,
            
            t.created_at, t.updated_at
        FROM tasks t
        WHERE t.id = $1
    `

    var t domain.Task
    err := r.db.QueryRow(ctx, query, id).Scan(
        &t.ID, &t.UserID, &t.Title, &t.Description, &t.Status,
        &t.Priority, &t.Labels, &t.ReminderTime,
        &t.RecurrencePattern,
        &t.NextRun,
        &t.CreatedAt, &t.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return &t, nil
}

// Update task (support update status, priority, labels, reminder_time, recurrence)
func (r *PostgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
    query := `
        UPDATE tasks 
        SET title = $1, description = $2, status = $3, priority = $4, labels = $5, reminder_time = $6, recurrence_pattern = $7, next_run = $8, updated_at = $9
        WHERE id = $10
    `

    cmdTag, err := r.db.Exec(ctx, query,
        task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.Labels,
        task.ReminderTime,
        task.RecurrencePattern,
        task.NextRun,
        task.UpdatedAt,
        task.ID,
    )
    if err != nil {
        return err
    }

    if cmdTag.RowsAffected() == 0 {
        return domain.ErrNotFound
    }

    return nil
}

// Delete task
func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM tasks WHERE id = $1`
    _, err := r.db.Exec(ctx, query, id)
    return err
}
