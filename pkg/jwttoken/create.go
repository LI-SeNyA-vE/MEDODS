package jwttoken

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(uuid, ip, accessSecret, refreshSecret string, lifetimeAccess, lifetimeRefresh time.Duration) (token *TokenDetails, err error) {
	td := &TokenDetails{}

	// Время жизни токенов
	atExpires := time.Now().Add(lifetimeAccess).Unix()  // Access живет ((time.Minute * 15) - 15 минут)
	rtExpires := time.Now().Add(lifetimeRefresh).Unix() // Refresh живет ((time.Hour * 24 * 1) - 1 день)

	// Генерация Access Token
	atClaims := jwt.MapClaims{
		"guid": uuid,
		"ip":   ip,
		"exp":  atExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	// Генерация Refresh Token
	rtClaims := jwt.MapClaims{
		"guid": uuid,
		"ip":   ip,
		"exp":  rtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}
