package users

import (
	"time"

	"github.com/google/uuid"
)

// ---------------------------------------------------------
// DATABASE ENTITIES
// ---------------------------------------------------------

// 1. User represents the core identity entity (Table: users)
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Omitted from JSON
	Department   string    `json:"department,omitempty"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Virtual field: Populated via SQL JOIN during login
	Roles []string `json:"roles,omitempty"`
}

// 2. Role represents a job function (Table: roles)
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 3. Permission represents an action on a resource (Table: permissions)
type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"` // e.g., "inventory:create"
	Description string `json:"description"`
}

// 4. UserRole represents the many-to-many mapping (Table: user_roles)
type UserRole struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID int       `json:"role_id"`
}

// 5. RolePermission represents the many-to-many mapping (Table: role_permissions)
type RolePermission struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}

// ---------------------------------------------------------
// HTTP REQUEST PAYLOADS
// ---------------------------------------------------------

// CreateUserRequest is the JSON payload received from the registration form
type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Department string `json:"department"`
}

// LoginRequest is the JSON payload received from the login form
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AssignRoleRequest is used by Admins to assign a new role to a user
type AssignRoleRequest struct {
	RoleID int `json:"role_id" validate:"required"`
}
