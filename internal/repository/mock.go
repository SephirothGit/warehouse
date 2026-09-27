package repository

import "context"

type MockUserRepo struct {
	CreateFunc          func(ctx context.Context, email, passwordHash string) (int, error)
	GetByEmailFunc      func(ctx context.Context, email string) (User, error)
	AssignRoleFunc      func(ctx context.Context, userID, roleID int) error
	GetRolesFunc        func(ctx context.Context, userID int) ([]string, error)
	GetRoleIDByNameFunc func(ctx context.Context, name string) (int, error)
}

func (m *MockUserRepo) Create(ctx context.Context, email, passwordHash string) (int, error) {
	return m.CreateFunc(ctx, email, passwordHash)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	return m.GetByEmailFunc(ctx, email)
}

func (m *MockUserRepo) AssignRole(ctx context.Context, userID, roleID int) error {
	return m.AssignRoleFunc(ctx, userID, roleID)
}

func (m *MockUserRepo) GetRoles(ctx context.Context, userID int) ([]string, error) {
	return m.GetRolesFunc(ctx, userID)
}

func (m *MockUserRepo) GetRoleIDByName(ctx context.Context, name string) (int, error) {
	return m.GetRoleIDByNameFunc(ctx, name)
}
