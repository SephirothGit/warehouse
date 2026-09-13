package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type UserRepo interface {
	Create(ctx context.Context, email, passwordHash string) (int, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	AssignRole(ctx context.Context, userID int, roleID int) error
	GetRoles(ctx context.Context, userID int) ([]string, error)
	GetRoleIDByName(ctx context.Context, name string) (int, error)
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, email string, passwordHash string) (int, error) {
	row := u.db.QueryRowContext(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, ErrAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, email, password_hash FROM users WHERE email = $1", email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (u *userRepo) AssignRole(ctx context.Context, userID int, roleID int) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) GetRoles(ctx context.Context, userID int) ([]string, error) {
	rows, err := u.db.QueryContext(ctx, "SELECT r.name FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var roleName string
		err := rows.Scan(&roleName)
		if err != nil {
			return nil, err
		}
		results = append(results, roleName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (u *userRepo) GetRoleIDByName(ctx context.Context, name string) (int, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id FROM roles WHERE name = $1", name)
	
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}