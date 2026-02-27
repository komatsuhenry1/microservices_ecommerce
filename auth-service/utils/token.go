package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
	"fmt"
	"log"
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

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// MODIFICAÇÃO IMPORTANTE AQUI!
	if err != nil {
		// Imprime o erro detalhado da biblioteca JWT no seu terminal
		log.Printf("Erro detalhado da validação do token: %v", err)
		// Agora podemos retornar a mensagem genérica para o frontend
		return nil, fmt.Errorf("Token inválido ou expirado")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Token inválido")
}