package main

import (
	serverconfig "MEDODS/internal/config/server"
	"MEDODS/internal/server/delivery/httpapi"
	"MEDODS/internal/server/repository/database"
	"MEDODS/internal/server/repository/database/postgresql"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	var storage *database.Storage
	var err error

	// Логгер
	l := logrus.New()
	log := logrus.NewEntry(l)
	log.Info("Запущен сервер")

	// Загрузка конфига
	cfg := serverconfig.NewConfigServerMock()

	//cfg, err := serverconfig.NewConfigServer()

	if err != nil {
		log.Fatal(err)
	}
	//Попытка подключения к БД PostgreSQL несколько раз.
	for i := 0; i < 3; i++ {
		//Создание переменной для работы с БД
		storage, err = postgresql.NewPostgresStorage(cfg.DatabaseCfg, log)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Проверка переменной подключения к БД
	if storage == nil {
		err = fmt.Errorf("не удалось подключиться к хранилищу: %w", err)
		log.Fatal(err)
		os.Exit(1)
	}

	defer storage.AuthToken.Close()

	log.Infof("Сервер запущен на %s", cfg.FlagAddressAndPort)
	err = httpapi.StartServerHTTP(storage, cfg, log)
	if err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
		return
	}
}
