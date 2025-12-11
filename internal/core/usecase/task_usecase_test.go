package usecase_test

import (
	"context"
	"testing"
	"time"

	"simple-task-manager/internal/core/domain"
	"simple-task-manager/internal/core/usecase"
	"simple-task-manager/internal/core/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	mockTaskRepo := new(mocks.TaskRepository)
	timeout := 2 * time.Second
	u := usecase.NewTaskUsecase(mockTaskRepo, timeout)

	t.Run("Success Create Task", func(t *testing.T) {
		task := &domain.Task{
			UserID: 1,
			Title:  "Belajar Golang",
		}

		// Expectation: Repo Create dipanggil sekali
		mockTaskRepo.On("Create", mock.Anything, task).Return(nil).Once()

		err := u.Create(context.Background(), task)

		assert.NoError(t, err)
		assert.Equal(t, "pending", task.Status) // Default status harus pending
	})
}