package authlogic_test

import (
	"MEDODS/internal/config/server"
	"MEDODS/internal/server/app/auth"
	"MEDODS/internal/server/repository/database"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
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

	// Создаём пару токенов
	tokens, err := authLogic.CreateToken(guid, ip)
	if err != nil {
		t.Fatalf("не удалось создать токен: %v", err)
	}

	time.Sleep(1 * time.Second)

	// Первый вызов — успех
	newTokens, err := authLogic.RefreshToken(tokens.RefreshToken, guid, ip, ip)
	if err != nil {
		t.Fatalf("первое обновление не удалось: %v", err)
	}

	if newTokens == nil || newTokens.RefreshToken == "" {
		t.Fatal("новый токен не передан")
	}

	time.Sleep(1 * time.Second)

	// Повторный вызов со старым токеном — должен упасть
	_, err = authLogic.RefreshToken(tokens.RefreshToken, guid, ip, ip)
	if err == nil {
		t.Fatal("токен прошёл валидацию, хотя не должен был")
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

	time.Sleep(1 * time.Second)

	// Новый IP → триггер срабатывает
	newTokens, err := authLogic.RefreshToken(tokens.RefreshToken, guid, ip, "10.10.10.10")
	if err != nil {
		t.Fatalf("Ошибка при рефреше токена: %v", err)
	}

	time.Sleep(1 * time.Second)

	// Старый Refresh
	_, err = authLogic.RefreshToken(tokens.RefreshToken, guid, ip, "10.10.10.10")
	if err == nil {
		t.Fatalf("Ошибка неволидный токен, прошёл валидацию: %v", err)
	}

	time.Sleep(1 * time.Second)

	_, err = authLogic.RefreshToken(newTokens.RefreshToken, guid, ip, "192.168.0.1")
	if err != nil {
		t.Fatalf("Ошибка при рефреше токена: %v", err)
	}

}
