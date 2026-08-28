package service

import "github.com/SephirothGit/warehouse/internal/repository"

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

