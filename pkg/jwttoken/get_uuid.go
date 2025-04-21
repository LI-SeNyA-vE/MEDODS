package jwttoken

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func GetClaim(token *jwt.Token, claim string) (uuid string, err error) {
	claims, _ := token.Claims.(jwt.MapClaims)
	uuid, ok := claims[claim].(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("%s не указан", claim))
	}
	return uuid, nil
}
