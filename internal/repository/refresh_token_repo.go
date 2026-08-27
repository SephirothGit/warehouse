package repository

import (
	"database/sql"
	"time"
)

type RefreshTokenRepo interface {
	Save(userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(tokenHash string) (userID int, expiresAt time.Time, revoked bool, err error)
	Revoke(tokenHash string) error
	RevokeAllForUser(userID int) error
}

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) RefreshTokenRepo {
	return &refreshTokenRepo{
		db: db,
	}
}

func (r *refreshTokenRepo) Save(userID int, tokenHash string, expiresAt time.Time) error {
_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)", userID, tokenHash, expiresAt)
if err != nil {
	return err
}
return nil
}

func (r *refreshTokenRepo) FindByHash(tokenHash string) (userID int, expiresAt time.Time, revoked bool, err error) {
	row := r.db.QueryRow("SELECT user_id, expires_at, revoked FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	err = row.Scan(&userID, &expiresAt, &revoked)
	return
}

func (r *refreshTokenRepo) Revoke(tokenHash string) error {
_, err := r.db.Exec("UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1", tokenHash)
if err != nil {
	return err
}
return nil
}

func (r *refreshTokenRepo) RevokeAllForUser(userID int) error {
_, err := r.db.Exec("UPDATE refresh_tokens SET revoked = true WHERE user_id = $1", userID)
if err != nil {
	return err
}
return nil
}