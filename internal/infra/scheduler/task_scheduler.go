package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"simple-task-manager/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskScheduler struct {
	db *pgxpool.Pool
}

func NewTaskScheduler(db *pgxpool.Pool) *TaskScheduler {
	return &TaskScheduler{db: db}
}

// ProcessRecurringTasks akan dipanggil setiap kali Cron berjalan
func (s *TaskScheduler) ProcessRecurringTasks() {
	ctx := context.Background()
	log.Println(">>> [CRON] Checking for recurring tasks...")

	// 1. Cari task yang recurrence pattern-nya aktif DAN waktunya sudah lewat (NextRun < Now)
	query := `
		SELECT id, user_id, title, description, priority, labels, recurrence_pattern, next_run
		FROM tasks 
		WHERE recurrence_pattern != '' 
		AND next_run <= NOW()
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Printf("[CRON ERROR] Fetching tasks: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t domain.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Priority, &t.Labels, &t.RecurrencePattern, &t.NextRun)
		if err != nil {
			continue
		}

		// 2. Clone Task Baru
		s.createNextInstance(ctx, &t)

		// 3. Update NextRun di Task Induk (Supaya tidak digenerate terus menerus)
		s.updateNextRun(ctx, &t)
	}
}

func (s *TaskScheduler) createNextInstance(ctx context.Context, parent *domain.Task) {
	// Buat task baru berdasarkan parent, tapi reset status jadi pending
	insertQuery := `
		INSERT INTO tasks (user_id, title, description, status, priority, labels, created_at, updated_at)
		VALUES ($1, $2, $3, 'pending', $4, $5, NOW(), NOW())
	`
	// Tambahkan penanda di judul (Opsional)
	newTitle := fmt.Sprintf("%s (Auto)", parent.Title)

	_, err := s.db.Exec(ctx, insertQuery, parent.UserID, newTitle, parent.Description, parent.Priority, parent.Labels)
	if err != nil {
		log.Printf("[CRON ERROR] Failed creating instance for task %d: %v", parent.ID, err)
	} else {
		log.Printf("[CRON SUCCESS] Created new instance for task: %s", parent.Title)
	}
}

func (s *TaskScheduler) updateNextRun(ctx context.Context, t *domain.Task) {
	var nextTime time.Time
	
	// Hitung waktu berikutnya berdasarkan pattern
	switch t.RecurrencePattern {
	case "daily":
		nextTime = t.NextRun.AddDate(0, 0, 1) // Tambah 1 hari
	case "weekly":
		nextTime = t.NextRun.AddDate(0, 0, 7) // Tambah 7 hari
	case "monthly":
		nextTime = t.NextRun.AddDate(0, 1, 0) // Tambah 1 bulan
	default:
		return // Pattern tidak dikenal
	}

	// Update DB
	_, err := s.db.Exec(ctx, "UPDATE tasks SET next_run = $1 WHERE id = $2", nextTime, t.ID)
	if err != nil {
		log.Printf("[CRON ERROR] Failed update next_run for task %d: %v", t.ID, err)
	}
}