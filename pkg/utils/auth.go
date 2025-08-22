// package utils

// import (
// 	"golang.org/x/crypto/bcrypt"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func Hash(pw string) ([]byte, error) { return bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost) }
// func CheckHash(pw string, hash []byte) error { return bcrypt.CompareHashAndPassword(hash, []byte(pw)) }

// import "github.com/golang-jwt/jwt/v5"
// func GetToken(uid int) (string, error){}
//

// pkg/utils/auth.go
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Version is set by ldflags (Makefile). Keep in this package so LDFLAGS can inject it.

// JWT secret used to sign tokens. Read from env; default is only for local/dev.
var jwtSecret = []byte(getEnv("JWT_SECRET", "dev-secret"))

// ErrInvalidToken returned when token parsing/validation fails.
var ErrInvalidToken = errors.New("invalid token")

// HashPassword returns a bcrypt hash of the given plaintext password.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success.
func CheckPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// GenerateToken creates a signed JWT containing the user ID and expiry.
// Expiration is 72 hours by default.
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"ver":     Version,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

// ParseToken validates a token and returns the user ID on success.
// It returns ErrInvalidToken when the token is invalid.
func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// ensure signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	// Note: jwt library unmarshals numbers to float64
	uidf, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return int(uidf), nil
}

// helper to read env with fallback
func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
