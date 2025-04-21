package database

type Storage struct {
	AuthToken AuthToken
}

type AuthToken interface {
	AddRefreshToken(guid string, refresh []byte) (err error)
	SearchRefreshToken(guid string) (refresh []byte, err error)
	DeleteRefreshToken(guid string) (err error)
	Close()
}
