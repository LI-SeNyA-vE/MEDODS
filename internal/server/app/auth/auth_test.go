package authlogic_test

import (
	"MEDODS/internal/config/server"
	"MEDODS/internal/server/app/auth"
	"MEDODS/internal/server/repository/database"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"testing"
)

type mockRepo struct {
	storage map[string][]byte
}

func (m *mockRepo) AddRefreshToken(guid string, token []byte) error {
	m.storage[guid] = token
	return nil
}

func (m *mockRepo) SearchRefreshToken(guid string) ([]byte, error) {
	if t, ok := m.storage[guid]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *mockRepo) DeleteRefreshToken(guid string) error {
	delete(m.storage, guid)
	return nil
}

func (m *mockRepo) Close() {
}

func TestRefreshToken_SuccessAndFailOnReuse(t *testing.T) {
	// Конфиг
	cfg := &serverconfig.Server{
		FlagAccessKey:       "ACCESS_SECRET",
		FlagRefreshKey:      "REFRESH_SECRET",
		FlagLifetimeAccess:  900000000000, // 15m
		FlagLifetimeRefresh: 86400000000000,
	}

	// Репо и логика
	mock := &mockRepo{storage: make(map[string][]byte)}
	log := logrus.New().WithField("test", "refresh")
	authLogic := authlogic.NewAuthLogic((*database.Storage)(&struct {
		AuthToken database.AuthToken
	}{AuthToken: mock}), cfg, log)

	guid := "test-guid"
	ip := "127.0.0.1"

	// 1. Сначала создаём пару токенов
	tokens, err := authLogic.CreateToken(guid, ip)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	// Раскодируем refresh из base64
	refreshRaw, err := base64.StdEncoding.DecodeString(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}

	// 2. Первый вызов — успех
	newTokens, err := authLogic.RefreshToken(string(refreshRaw), guid, ip, ip)
	if err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}
	if newTokens == nil || newTokens.RefreshToken == "" {
		t.Fatal("new tokens not returned")
	}

	// 3. Повторный вызов с тем же токеном — должен упасть
	_, err = authLogic.RefreshToken(string(refreshRaw), guid, ip, ip)
	if err == nil {
		t.Fatal("second use of refresh token should fail, but passed")
	}
}

func TestRefreshToken_IPMismatch(t *testing.T) {
	cfg := &serverconfig.Server{
		FlagAccessKey:       "ACCESS_SECRET",
		FlagRefreshKey:      "REFRESH_SECRET",
		FlagLifetimeAccess:  900000000000,
		FlagLifetimeRefresh: 86400000000000,
	}

	mock := &mockRepo{storage: make(map[string][]byte)}
	log := logrus.New().WithField("test", "ip-mismatch")
	authLogic := authlogic.NewAuthLogic((*database.Storage)(&struct {
		AuthToken database.AuthToken
	}{AuthToken: mock}), cfg, log)

	guid := "test-guid"
	ip := "192.168.0.1"

	// Создаём токены
	tokens, err := authLogic.CreateToken(guid, ip)
	if err != nil {
		t.Fatalf("Ошибка создания токена: %v", err)
	}

	// Новый IP → триггер срабатывает
	newTokens, err := authLogic.RefreshToken(tokens.RefreshToken, guid, ip, "10.10.10.10")
	if err != nil {
		t.Fatalf("Ошибка при рефреше токена: %v", err)
	}

	// Старый Refresh
	_, err = authLogic.RefreshToken(tokens.RefreshToken, guid, ip, "10.10.10.10")
	if err == nil {
		t.Fatalf("Ошибка неволидный токен, прошёл валидацию: %v", err)
	}

	_, err = authLogic.RefreshToken(newTokens.RefreshToken, guid, ip, "192.168.0.1")
	if err != nil {
		t.Fatalf("Ошибка при рефреше токена: %v", err)
	}

}
