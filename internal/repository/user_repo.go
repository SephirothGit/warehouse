package repository

import (
	"database/sql"
)

type UserRepo interface {
	Create(email, passwordHash string) (int, error)
	GetByEmail(email string) (User, error)
	AssignRole(userID int, roleID int) error
	GetRoles(userID int) ([]string, error)
	GetRoleIDByName(name string) (int, error)
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

func (u *userRepo) Create(email string, passwordHash string) (int, error) {
	row := u.db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *userRepo) GetByEmail(email string) (User, error) {
	row := u.db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *userRepo) AssignRole(userID int, roleID int) error {
	_, err := u.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) GetRoles(userID int) ([]string, error) {
	rows, err := u.db.Query("SELECT r.name FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = $1", userID)
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

func (u *userRepo) GetRoleIDByName(name string) (int, error) {
	row := u.db.QueryRow("SELECT id FROM roles WHERE name = $1", name)
	
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}