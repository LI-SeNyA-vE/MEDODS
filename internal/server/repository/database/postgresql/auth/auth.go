package authDB

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func NewAuthDB(db *sql.DB, log *logrus.Entry) *AuthDB {
	return &AuthDB{
		db:  db,
		log: log,
	}
}

type AuthDB struct {
	db  *sql.DB
	log *logrus.Entry
}

func (a *AuthDB) AddRefreshToken(guid string, refresh []byte) (err error) {
	_, err = a.db.Exec(addRefreshToken, guid, refresh)
	a.log.Infof("Добавлен для GUID %s Refresh %s", guid, refresh)
	return err
}

func (a *AuthDB) SearchRefreshToken(guid string) (refresh []byte, err error) {
	err = a.db.QueryRow(querySearchRefresh, guid).Scan(&refresh)
	return refresh, err
}

func (a *AuthDB) DeleteRefreshToken(guid string) (err error) {
	res, err := a.db.Exec(queryDeleteRefresh, guid)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	a.log.Infof("Удалён %d Refresh %s", ra, guid)
	return err
}

func (a *AuthDB) Close() {
	a.db.Close()
}
