package middleware

import (
	serverconfig "MEDODS/internal/config/server"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	cfg *serverconfig.Server
	log *logrus.Entry
}

func NewMiddleware(cfg *serverconfig.Server, log *logrus.Entry) *Middleware {
	return &Middleware{
		cfg: cfg,
		log: log,
	}
}
