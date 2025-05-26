package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Anand078/rbac/internal/database"
	"github.com/Anand078/rbac/internal/models"
)

type AuthService struct {
	db        *database.DB
	jwtSecret string
}

func NewAuthService(db *database.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req models.CreateUserRequest) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	query := `
        INSERT INTO users (id, email, name, password_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING created_at, updated_at
    `
	err = tx.QueryRow(query, user.ID, user.Email, user.Name, user.PasswordHash).
		Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role if provided
	if req.RoleID != "" {
		roleID, err := uuid.Parse(req.RoleID)
		if err != nil {
			return nil, fmt.Errorf("invalid role ID: %w", err)
		}

		_, err = tx.Exec(
			"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
			user.ID, roleID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to assign role: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	query := `
        SELECT id, email, name, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	err := s.db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Load user roles
	roles, err := s.getUserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	// Generate JWT token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) getUserRoles(userID uuid.UUID) ([]models.Role, error) {
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
