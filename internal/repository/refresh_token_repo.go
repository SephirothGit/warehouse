package repository

import (
	"context"
	"database/sql"
	"time"
)

type RefreshTokenRepo interface {
	Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, revoked bool, err error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID int) error
}

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) RefreshTokenRepo {
	return &refreshTokenRepo{
		db: db,
	}
}

func (r *refreshTokenRepo) Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)", userID, tokenHash, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenRepo) FindByHash(ctx context.Context, tokenHash string) (userID int, expiresAt time.Time, revoked bool, err error) {
	row := r.db.QueryRowContext(ctx, "SELECT user_id, expires_at, revoked FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	err = row.Scan(&userID, &expiresAt, &revoked)
	return
}

func (r *refreshTokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1", tokenHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenRepo) RevokeAllForUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked = true WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}
