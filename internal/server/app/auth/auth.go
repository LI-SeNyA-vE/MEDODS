package authlogic

import (
	serverconfig "MEDODS/internal/config/server"
	"MEDODS/internal/server/repository/database"
	"MEDODS/pkg/jwttoken"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogic interface {
	CreateToken(guid string, ip string) (token *jwttoken.TokenDetails, err error)
	RefreshToken(refreshToken string, guid string, oldIP string, nowIP string) (token *jwttoken.TokenDetails, err error)
}

type authLogic struct {
	repo      database.AuthToken
	cfgServer *serverconfig.Server
	log       *logrus.Entry
}

func (a *authLogic) CreateToken(guid, ip string) (token *jwttoken.TokenDetails, err error) {
	newToken, err := jwttoken.NewToken(guid,
		ip,
		a.cfgServer.FlagAccessKey,
		a.cfgServer.FlagRefreshKey,
		a.cfgServer.FlagLifetimeAccess,
		a.cfgServer.FlagLifetimeRefresh)
	if err != nil {
		return nil, err
	}

	ref, err := a.repo.SearchRefreshToken(guid)
	if ref != nil {
		err = a.repo.DeleteRefreshToken(guid)
		if err != nil {
			return nil, err
		}
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	//Для преодоления ограничений bcrypt в 72 байта, хешируем SHA-256
	shaHash := sha256.Sum256([]byte(newToken.RefreshToken))
	hexHash := hex.EncodeToString(shaHash[:])

	//Создаём хеш для сохранения в БД
	hashed, err := bcrypt.GenerateFromPassword([]byte(hexHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = a.repo.AddRefreshToken(guid, hashed)
	if err != nil {
		return nil, err
	}

	//Перезаписываем newToken кодирую Refresh
	newToken.RefreshToken = base64.StdEncoding.EncodeToString([]byte(newToken.RefreshToken))

	return newToken, nil
}

func (a *authLogic) RefreshToken(refreshToken string, guid string, oldIP string, nowIP string) (token *jwttoken.TokenDetails, err error) {
	if oldIP != nowIP {
		a.sendIPAlert(guid, oldIP, nowIP)
	}

	dbRefreshToken, err := a.repo.SearchRefreshToken(guid)
	if err != nil {
		return nil, err
	}

	a.log.Infof("Переданный Refresh: %s", refreshToken)
	//Для преодоления ограничений bcrypt в 72 байта, хешируем SHA-256

	//Раскадируем из base64
	byteRefreshToken, _ := base64.StdEncoding.DecodeString(refreshToken)
	shaHash := sha256.Sum256(byteRefreshToken)
	hexHash := hex.EncodeToString(shaHash[:])

	a.log.Infof("Закодированный переданный Refresh: %s", hexHash)
	a.log.Infof("Hash сохранённого в БД Refresh: %s", dbRefreshToken)

	err = bcrypt.CompareHashAndPassword(dbRefreshToken, []byte(hexHash))
	if err != nil {
		return nil, errors.New("refresh token не последний созданный")
	}

	a.log.Info("Токены равны")

	newToken, err := jwttoken.NewToken(guid,
		nowIP,
		a.cfgServer.FlagAccessKey,
		a.cfgServer.FlagRefreshKey,
		a.cfgServer.FlagLifetimeAccess,
		a.cfgServer.FlagLifetimeRefresh)
	if err != nil {
		return nil, err
	}

	err = a.repo.DeleteRefreshToken(guid)
	if err != nil {
		return nil, err
	}

	a.log.Infof("Создан новый Refresh: %s \nУдалён старый", newToken.RefreshToken)

	//Для преодоления ограничений bcrypt в 72 байта, хешируем SHA-256
	shaHash = sha256.Sum256([]byte(newToken.RefreshToken))
	hexHash = hex.EncodeToString(shaHash[:])

	a.log.Infof("Закодированный созданный Refresh: %s", hexHash)
	//Создаём хеш для сохранения в БД
	hashed, err := bcrypt.GenerateFromPassword([]byte(hexHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	a.log.Infof("Hash который сохраним в БД Refresh: %s", hashed)

	err = a.repo.AddRefreshToken(guid, hashed)
	if err != nil {
		return nil, err
	}

	//Перезаписываем newToken кодируя Refresh
	newToken.RefreshToken = base64.StdEncoding.EncodeToString([]byte(newToken.RefreshToken))

	a.log.Infof("Отправили пользователю Refresh: %s \n\n\n", newToken.RefreshToken)

	return newToken, nil
}

func (a *authLogic) sendIPAlert(guid, oldIP, newIP string) {
	// тут просто логируем, но можно вставить SMTP‑реализацию
	a.log.Infof("[ВНИМАНИЕ] у guid=%s изменился IP с %s на %s", guid, oldIP, newIP)
}

func NewAuthLogic(repo *database.Storage, cfgServer *serverconfig.Server, log *logrus.Entry) AuthLogic {
	return &authLogic{
		repo:      repo.AuthToken,
		cfgServer: cfgServer,
		log:       log,
	}
}
