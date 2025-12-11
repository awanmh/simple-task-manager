package usecase_test

import (
	"context"
	// "errors" <-- HAPUS BARIS INI
	"testing"
	"time"

	"simple-task-manager/internal/core/domain"
	"simple-task-manager/internal/core/usecase"
	"simple-task-manager/internal/core/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	timeout := 2 * time.Second
	u := usecase.NewUserUsecase(mockUserRepo, timeout, "secret_key")

	t.Run("Success Register", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(nil, nil).Once()
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		err := u.Register(context.Background(), user)

		assert.NoError(t, err)
		assert.NotEqual(t, "password123", user.Password) 
		
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failed - Email Already Exists", func(t *testing.T) {
		existingUser := &domain.User{ID: 1, Email: "test@example.com"}
		mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(existingUser, nil).Once()

		err := u.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
		
		mockUserRepo.AssertNotCalled(t, "Create")
	})
}