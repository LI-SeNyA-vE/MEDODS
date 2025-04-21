package postgresql

import (
	"database/sql"
	"fmt"
)

func createTableIsNot(db *sql.DB) (err error) {
	query, err := db.Query(queryTokenTableIsNot)
	if err != nil {
		return fmt.Errorf("ошибка при создание таблицы: %w", err)
	}
	query.Close()
	return nil
}
