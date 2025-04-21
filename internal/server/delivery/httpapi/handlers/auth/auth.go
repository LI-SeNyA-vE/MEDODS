package authhandlers

import (
	"MEDODS/internal/domain"
	"MEDODS/internal/server/app"
	"MEDODS/pkg/jwttoken"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type AuthHandler interface {
	PostAuthToken(w http.ResponseWriter, r *http.Request)
	PutRefresh(w http.ResponseWriter, r *http.Request)
}

// Реализация
type authHandler struct {
	biz app.BizLogic
	log *logrus.Entry
}

func (a authHandler) PostAuthToken(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var err error
	var resp domain.ResponseAuthToken

	if r.Header.Get("Content-Type") != "application/json" {
		err = errors.New("Сontent-Type не равен application/json")
		a.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("ошибка чтения тела запроса: %v", err)
		a.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		err = fmt.Errorf("ошибка разбора тела запроса в структуру: %v", err)
		a.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nowIP := r.RemoteAddr
	if strings.Contains(nowIP, ":") {
		nowIP = strings.Split(nowIP, ":")[0]
	}

	//Переходим на следующий слой логики
	token, err := a.biz.Auth.CreateToken(resp.GUID, nowIP)
	if err != nil {
		a.log.Errorf("кидаем ответ с ошибкой \"%v\"", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,                           // ❗ Доступен только серверу (защита от XSS)
		Secure:   false,                          // ❗ Только по HTTPS (обязательно в проде (true))
		SameSite: http.SameSiteStrictMode,        // ❗ Защита от CSRF
		Path:     "/api/refresh",                 // ❗ Ограничиваем путь (только для refresh)
		Expires:  time.Now().Add(24 * time.Hour), // ❗ Refresh-токен живёт 1 день
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы получили пару токенов, refresh записан в базу в закодированном виде base64"))

	//Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
	//Если пользователь повторно получает пару Access, Refresh токенов по тому же GUID, то мы удаляем старый и создаём новую
}

func (a authHandler) PutRefresh(w http.ResponseWriter, r *http.Request) {
	var err error

	cookie, _ := r.Cookie("refreshToken")

	validToken := r.Context().Value("validToken").(*jwt.Token)

	guid, err := jwttoken.GetClaim(validToken, "guid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	oldIP, err := jwttoken.GetClaim(validToken, "ip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nowIP := r.RemoteAddr
	if strings.Contains(nowIP, ":") {
		nowIP = strings.Split(nowIP, ":")[0]
	}

	token, err := a.biz.Auth.RefreshToken(cookie.Value, guid, oldIP, nowIP)
	if err != nil {
		a.log.Errorf("кидаем ответ с ошибкой \"%v\"", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,                           // ❗ Доступен только серверу (защита от XSS)
		Secure:   false,                          // ❗ Только по HTTPS (обязательно в проде (true))
		SameSite: http.SameSiteStrictMode,        // ❗ Защита от CSRF
		Path:     "/api/refresh",                 // ❗ Ограничиваем путь (только для refresh)
		Expires:  time.Now().Add(24 * time.Hour), // ❗ Refresh-токен живёт 1 день
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы рефрешнули и получили пару токенов, refresh записан в базу в закодированном виде base64"))

	//Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
}

func NewHandlers(uc app.BizLogic, log *logrus.Entry) AuthHandler {
	return &authHandler{
		biz: uc,
		log: log,
	}
}
