package postgresql

import (
	serverconfig "MEDODS/internal/config/server"
	"MEDODS/internal/server/repository/database"
	authDB "MEDODS/internal/server/repository/database/postgresql/auth"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func NewPostgresStorage(dbConnect serverconfig.DatabaseConfig, log *logrus.Entry) (*database.Storage, error) {
	// Подключаемся к дефолтной базе
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConnect.FlagDatabaseHost,
		dbConnect.FlagDatabasePort,
		dbConnect.FlagDatabaseLogin,
		dbConnect.FlagDatabasePassword,
		"postgres",
		dbConnect.FlagDatabaseSSLMode)) //"postgresql://Senya:1q2w3e4r5t@localhost:5432/postgres?sslmode=disable"
	if err != nil {
		err = fmt.Errorf("ошибка подключения к системной базе данных: %v", err)
		return nil, err
	}

	// Проверка соединения (ping)
	if err = db.Ping(); err != nil {
		err = fmt.Errorf("пинги в дефолтную базу не прошли: %w", err)
		return nil, err
	}

	// Проверяем существование базы
	var dbExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbConnect.FlagDatabaseName).Scan(&dbExists)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки базы: %w", err)
	}

	// Если базы нет - создаём
	if !dbExists {
		safeDBName := pgx.Identifier{dbConnect.FlagDatabaseName}.Sanitize()
		query := fmt.Sprintf("CREATE DATABASE %s", safeDBName)
		_, err = db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания базы: %w", err)
		}
	}

	//Подключаемся к нужной нам базе, которая точно есть
	db, err = sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConnect.FlagDatabaseHost,
		dbConnect.FlagDatabasePort,
		dbConnect.FlagDatabaseLogin,
		dbConnect.FlagDatabasePassword,
		dbConnect.FlagDatabaseName,
		dbConnect.FlagDatabaseSSLMode)) //"postgresql://Senya:1q2w3e4r5t@localhost:5432/MEDODS?sslmode=disable"
	if err != nil {
		err = fmt.Errorf("ошибка подключения к системной базе данных: %v", err)
		return nil, err
	}

	// Проверка соединения (ping)
	if err = db.Ping(); err != nil {
		err = fmt.Errorf("пинги в нашу базу не прошли: %w", err)
		return nil, err
	}

	err = createTableIsNot(db)
	if err != nil {
		return nil, err
	}

	return &database.Storage{
		AuthToken: authDB.NewAuthDB(db, log),
	}, nil
}
