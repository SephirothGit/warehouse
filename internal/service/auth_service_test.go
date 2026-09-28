package service

import (
	"context"
	"testing"

	"github.com/SephirothGit/warehouse/internal/repository"
)

func TestRegister_Success(t *testing.T) {
	mockUserRepo := &repository.MockUserRepo{
		CreateFunc: func(ctx context.Context, email, passwordHash string) (int, error) {
			return 1, nil
		},
		GetRoleIDByNameFunc: func(ctx context.Context, name string) (int, error) {
			return 3, nil
		},
		AssignRoleFunc: func(ctx context.Context, userID, roleID int) error {
			return nil
		},
	}

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
	mockUserRepo := repository.MockUserRepo{
		CreateFunc: func(ctx context.Context, email, passwordHash string) (int, error) {
			return 0, repository.ErrAlreadyExists
		},
	}

	authService := NewAuthService(mockUserRepo, nil, []byte("test-secret"))

	_, err := authService.Register(context.Background(), "test@test.com", "password123")
	if err != nil {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}
