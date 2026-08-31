package service

import (
	"fmt"
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
	RefreshAccessToken(rawRefreshToken string) (TokenPair, error)
	Logout(rawRefreshToken string) error
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

func (a *authService) RefreshAccessToken(rawRefreshToken string) (TokenPair, error) {
	tokenHash := auth.HashToken(rawRefreshToken)

	userID, expiresAt, revoked, err := a.refreshRepo.FindByHash(tokenHash)
	if err != nil {
		return TokenPair{}, err
	}

	if revoked {
		a.refreshRepo.RevokeAllForUser(userID)
		return TokenPair{}, fmt.Errorf("token reuse detected, all sessions revoked")
	}

	if expiresAt.Before(time.Now()) {
		return TokenPair{}, fmt.Errorf("refresh token expired")
	}

	err = a.refreshRepo.Revoke(tokenHash)
	if err != nil {
		return TokenPair{}, err
	}

	newRawRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	newTokenHash := auth.HashToken(newRawRefreshToken)
	newExpiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = a.refreshRepo.Save(userID, newTokenHash, newExpiresAt)
	if err != nil {
		return TokenPair{}, err
	}

	accessToken, err := auth.GenerateToken(strconv.Itoa(userID), a.jwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: newRawRefreshToken}, nil
}

func (a *authService) Logout(rawRefreshToken string) error {
	tokenHash := auth.HashToken(rawRefreshToken)
	return a.refreshRepo.Revoke(tokenHash)
}
