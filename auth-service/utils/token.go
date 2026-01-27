package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateToken(userId string, userRole string, userName string, expiration time.Duration) (string, error) {
	expiresAt := time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    userId,
		"role":   userRole,
		"name":   userName,
		"exp":    expiresAt,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
