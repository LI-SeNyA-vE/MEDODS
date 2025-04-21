package httpapi

import (
	serverconfig "MEDODS/internal/config/server"
	"MEDODS/internal/server/app"
	"MEDODS/internal/server/delivery/httpapi/handlers"
	"MEDODS/internal/server/delivery/httpapi/middleware"
	"MEDODS/internal/server/delivery/httpapi/router"
	"MEDODS/internal/server/repository/database"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartServerHTTP(storage *database.Storage, cfg *serverconfig.Server, log *logrus.Entry) (err error) {
	uc := app.NewBizLogic(storage, cfg, log)
	handle := handlers.NewHandlers(uc, log)
	mw := middleware.NewMiddleware(cfg, log)
	r := router.NewRouter(log, mw, handle)
	r.SetupRouter()

	err = http.ListenAndServe(cfg.FlagAddressAndPort, r.Mux)
	if err != nil {
		return err
	}
	return nil
}
