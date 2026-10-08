package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/SephirothGit/warehouse/internal/mocks"
	"github.com/SephirothGit/warehouse/internal/repository"
)

func TestRegister_Success(t *testing.T) {
	mockUserRepo := mocks.NewUserRepo(t)

	mockUserRepo.EXPECT().
		Create(mock.Anything, "test@test.com", mock.Anything).
		Return(1, nil)

	mockUserRepo.EXPECT().
		GetRoleIDByName(mock.Anything, "warehouse_worker").
		Return(3, nil)

	mockUserRepo.EXPECT().
		AssignRole(mock.Anything, 1, 3).
		Return(nil)

	authService := NewAuthService(mockUserRepo, nil, []byte("test-secret"))

	userID, err := authService.Register(context.Background(), "test@test.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if userID != 1 {
		t.Errorf("expected userID 1, got %d", userID)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := mocks.NewUserRepo(t)

	mockUserRepo.EXPECT().
		Create(mock.Anything, "test@test.com", mock.Anything).
		Return(0, repository.ErrAlreadyExists)

	authService := NewAuthService(mockUserRepo, nil, []byte("test-secret"))

	_, err := authService.Register(context.Background(), "test@test.com", "password123")
	if !errors.Is(err, repository.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}
