package jwttoken

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken Функция проверки токена
func ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи - SHA512
		if token.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("ожидался другой метод подписи")
		}
		return []byte(secret), nil
	})

	// Проверяем сам токен на валидность
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("токен неволиден, возможно отредактирован/подменён: %w", NoValid)
	}

	return token, nil
}
