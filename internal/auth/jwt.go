package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Claims defines the payload of our JWT
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a token for a user with their roles
func GenerateJWT(userID uuid.UUID, roles []string) (string, error) {
	if len(jwtKey) == 0 {
		jwtKey = []byte("my-secret-key") // Fallback for dev
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyJWT parses and validates a JWT string
func VerifyJWT(tokenString string) (*Claims, error) {
	if len(jwtKey) == 0 {
		jwtKey = []byte("my-secret-key") // Fallback for dev
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
