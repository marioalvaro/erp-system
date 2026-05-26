// internal/users/repository.go
package users

import (
	"database/sql"
	"errors"
)

// GetRoleIDByName queries PostgreSQL to find the dynamic ID of a specific role
func GetRoleIDByName(db *sql.DB, roleName string) (int, error) {
	var id int
	
	// Query the roles table for the ID matching the string name
	query := `SELECT id FROM roles WHERE name = $1`
	err := db.QueryRow(query, roleName).Scan(&id)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("role not found")
		}
		return 0, err
	}
	
	return id, nil
}

// AssignRoleToUser maps a user to a specific role in the user_roles table
func AssignRoleToUser(db *sql.DB, userID interface{}, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Exec(query, userID, roleID)
	return err
}

// CreateUser inserts a new user into the PostgreSQL database
func CreateUser(db *sql.DB, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, department)
		VALUES ($1, $2, $3)
		RETURNING id, is_verified, created_at, updated_at
	`
	
	err := db.QueryRow(
		query, 
		user.Email, 
		user.PasswordHash, 
		user.Department,
	).Scan(&user.ID, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Error handling for duplicate emails
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.New("a user with this email already exists")
		}
		return err
	}

	return nil
}