package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}
