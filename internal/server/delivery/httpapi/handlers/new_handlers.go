package handlers

import (
	"MEDODS/internal/server/app"
	"MEDODS/internal/server/delivery/httpapi/handlers/auth"
	"github.com/sirupsen/logrus"
)

// Handlers хранит ссылки на логгер (logrus.Entry)
// и реализацию интерфейса UserRepository (repository).
// Используется в хендлерах для взаимодействия
// с базой, а также для логирования запросов/ответов.
type Handlers struct {
	Auth authhandlers.AuthHandler
}

func NewHandlers(uc app.BizLogic, log *logrus.Entry) *Handlers {
	return &Handlers{
		Auth: authhandlers.NewHandlers(uc, log),
	}
}
