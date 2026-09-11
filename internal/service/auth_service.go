package service

import (
	"context"
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
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email, password string) (TokenPair, error)
	RefreshAccessToken(ctx context.Context, rawRefreshToken string) (TokenPair, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

func NewAuthService(userRepo repository.UserRepo, refreshRepo repository.RefreshTokenRepo, jwtSecret []byte) AuthService {
	return &authService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtSecret:   jwtSecret,
	}
}

func (a *authService) Register(ctx context.Context, email, password string) (int, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	userID, err := a.userRepo.Create(ctx, email, hashedPassword)
	if err != nil {
		return 0, err
	}

	roleID, err := a.userRepo.GetRoleIDByName(ctx, "warehouse_worker")
	if err != nil {
		return 0, err
	}

	err = a.userRepo.AssignRole(ctx, userID, roleID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (a *authService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
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

	err = a.refreshRepo.Save(ctx, user.ID, tokenHash, expiresAt)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *authService) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (TokenPair, error) {
	tokenHash := auth.HashToken(rawRefreshToken)

	userID, expiresAt, revoked, err := a.refreshRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return TokenPair{}, err
	}

	if revoked {
		a.refreshRepo.RevokeAllForUser(ctx, userID)
		return TokenPair{}, fmt.Errorf("token reuse detected, all sessions revoked")
	}

	if expiresAt.Before(time.Now()) {
		return TokenPair{}, fmt.Errorf("refresh token expired")
	}

	err = a.refreshRepo.Revoke(ctx, tokenHash)
	if err != nil {
		return TokenPair{}, err
	}

	newRawRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	newTokenHash := auth.HashToken(newRawRefreshToken)
	newExpiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = a.refreshRepo.Save(ctx, userID, newTokenHash, newExpiresAt)
	if err != nil {
		return TokenPair{}, err
	}

	accessToken, err := auth.GenerateToken(strconv.Itoa(userID), a.jwtSecret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: newRawRefreshToken}, nil
}

func (a *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := auth.HashToken(rawRefreshToken)
	return a.refreshRepo.Revoke(ctx, tokenHash)
}
