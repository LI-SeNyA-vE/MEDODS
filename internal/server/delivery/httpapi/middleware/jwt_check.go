package middleware

import (
	"MEDODS/pkg/jwttoken"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

func (m *Middleware) JwtCheck(h http.Handler) http.Handler {
	jwtCheck := func(w http.ResponseWriter, r *http.Request) {
		baseRefreshToken, _ := r.Cookie("refreshToken")
		if baseRefreshToken == nil || baseRefreshToken.Value == "" {
			http.Error(w, "переданы невальдные куки", http.StatusUnauthorized)
			return
		}

		//Раскадируем из base64
		refreshToken, _ := base64.StdEncoding.DecodeString(baseRefreshToken.Value)

		// Проверяем токен на валидность
		validToken, err := jwttoken.ValidateToken(string(refreshToken), m.cfg.FlagRefreshKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("токен не прошёл валидацию: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "validToken", validToken)

		// Передаём управление дальше
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(jwtCheck)
}
