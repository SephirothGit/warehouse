package service

import (
	"strconv"
	"time"

	"github.com/SephirothGit/warehouse/internal/auth"
	"github.com/SephirothGit/warehouse/internal/repository"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type authService struct {
	userRepo    repository.UserRepo
	refreshRepo repository.RefreshTokenRepo
	jwtSecret   []byte
}

type AuthService interface {
	Register(email, password string) (int, error)
	Login(email, password string) (TokenPair, error)
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

func (a *authService) Login(email, password string) (TokenPair, error) {
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return TokenPair{}, err
	}

	err = auth.CheckPassword(user.PasswordHash, password)
	if err != nil {
		return TokenPair{}, err
	}

	accessToken, err := auth.GenerateToken(strconv.Itoa(user.ID), a.jwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	tokenHash := auth.HashToken(refreshToken)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = a.refreshRepo.Save(user.ID, tokenHash, expiresAt)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
