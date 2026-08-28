package service

import (
	"github.com/SephirothGit/warehouse/internal/auth"
	"github.com/SephirothGit/warehouse/internal/repository"
)

type authService struct {
	userRepo    repository.UserRepo
	refreshRepo repository.RefreshTokenRepo
	jwtSecret   []byte
}

type AuthService interface {
	Register(email, password string) (int, error)
}

func NewAuthService(userRepo repository.UserRepo, refreshRepo repository.RefreshTokenRepo, jwtSecret []byte) AuthService {
	return &authService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtSecret:   jwtSecret,
	}
}

func (a *authService) Register(email, password string) (int, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	userID, err := a.userRepo.Create(email, hashedPassword)
	if err != nil {
		return 0, err
	}

	roleID, err := a.userRepo.GetRoleIDByName("warehouse_worker")
	if err != nil {
		return 0, err
	}

	err = a.userRepo.AssignRole(userID, roleID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}