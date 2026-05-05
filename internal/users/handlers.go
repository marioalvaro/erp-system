package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
	"github.com/marioalvaro/erp-system/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sql.DB
}

// RegisterUser handles POST /api/v1/auth/register
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	newUser := User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Department:   req.Department,
	}

	if err := CreateUser(h.DB, &newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful. Please verify your email.",
	})
}

// LoginUser handles POST /api/v1/auth/login
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user User
	query := `
		SELECT u.id, u.email, u.password_hash, u.is_verified, array_agg(r.name) 
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.email = $1
		GROUP BY u.id`

	var roles pq.StringArray
	err := h.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsVerified, &roles)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	user.Roles = []string(roles)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !user.IsVerified {
		http.Error(w, "Please verify your email before logging in", http.StatusForbidden)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Roles)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"message": "Login successful",
	})
}

// GetProfile handles GET /api/v1/users/profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch full user details if needed, but for now just return claims
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claims)
}
