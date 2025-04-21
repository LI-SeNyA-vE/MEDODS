package app

import (
	serverconfig "MEDODS/internal/config/server"
	authlogic "MEDODS/internal/server/app/auth"
	"MEDODS/internal/server/repository/database"
	"github.com/sirupsen/logrus"
)

type BizLogic struct {
	Auth authlogic.AuthLogic
}

func NewBizLogic(repo *database.Storage, cfgServer *serverconfig.Server, log *logrus.Entry) BizLogic {
	return BizLogic{
		Auth: authlogic.NewAuthLogic(repo, cfgServer, log),
	}
}
