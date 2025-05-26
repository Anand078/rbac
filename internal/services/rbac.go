package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/Anand078/rbac/internal/database"
	"github.com/Anand078/rbac/internal/models"
)

type RBACService struct {
	db *database.DB
}

func NewRBACService(db *database.DB) *RBACService {
	return &RBACService{db: db}
}

// Role Management
func (s *RBACService) CreateRole(req models.CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	query := `
        INSERT INTO roles (id, name, description)
        VALUES ($1, $2, $3)
        RETURNING created_at
    `
	err := s.db.QueryRow(query, role.ID, role.Name, role.Description).Scan(&role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (s *RBACService) GetAllRoles() ([]models.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles ORDER BY name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *RBACService) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	query := `
        SELECT r.id, r.name, r.description, r.created_at
        FROM roles r
        JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1
    `
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *RBACService) AssignRole(userID, roleID uuid.UUID) error {
	query := `
        INSERT INTO user_roles (user_id, role_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, role_id) DO NOTHING
    `
	_, err := s.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}
	return nil
}

func (s *RBACService) RemoveRole(userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := s.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	return nil
}

// Permission Management
func (s *RBACService) CreatePermission(req models.CreatePermissionRequest) (*models.Permission, error) {
	permission := &models.Permission{
		ID:          uuid.New(),
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	query := `
        INSERT INTO permissions (id, name, resource, action, description)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING created_at
    `
	err := s.db.QueryRow(query, permission.ID, permission.Name, permission.Resource,
		permission.Action, permission.Description).Scan(&permission.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

func (s *RBACService) GetAllPermissions() ([]models.Permission, error) {
	query := `SELECT id, name, resource, action, description, created_at FROM permissions ORDER BY resource, action`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action,
			&perm.Description, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *RBACService) GetRolePermissions(roleID uuid.UUID) ([]models.Permission, error) {
	query := `
        SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
        ORDER BY p.resource, p.action
    `
	rows, err := s.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action,
			&perm.Description, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *RBACService) GrantPermission(roleID, permissionID uuid.UUID) error {
	query := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
        ON CONFLICT (role_id, permission_id) DO NOTHING
    `
	_, err := s.db.Exec(query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}
	return nil
}

func (s *RBACService) RevokePermission(roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := s.db.Exec(query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}
	return nil
}

// Authorization Check
func (s *RBACService) HasPermission(userID uuid.UUID, resource, action string) (bool, error) {
	query := `
        SELECT COUNT(*) > 0
        FROM users u
        JOIN user_roles ur ON u.id = ur.user_id
        JOIN role_permissions rp ON ur.role_id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = $1 AND p.resource = $2 AND p.action = $3
    `
	var hasPermission bool
	err := s.db.QueryRow(query, userID, resource, action).Scan(&hasPermission)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return hasPermission, nil
}
