package serverconfig

import "time"

type Server struct {
	FlagAddressAndPort  string         `envconfig:"ADDRESS" default:"localhost:8901"`
	FlagLogLevel        string         `envconfig:"LOG_LEVEL" default:"info"`
	DatabaseCfg         DatabaseConfig `json:"database_cfg"`
	FlagAccessKey       string         `envconfig:"ACCESS_KEY" required:"true"`
	FlagLifetimeAccess  time.Duration  `envconfig:"LIFETIME_ACCESS" default:"15m"`
	FlagRefreshKey      string         `envconfig:"REFRESH_KEY" required:"true"`
	FlagLifetimeRefresh time.Duration  `envconfig:"LIFETIME_REFRESH" default:"24h"`
}

type DatabaseConfig struct {
	FlagDatabaseHost     string `envconfig:"DATABASE_HOST" json:"database_host" default:"localhost"`
	FlagDatabasePort     string `envconfig:"DATABASE_PORT" json:"database_port" default:"5432"`
	FlagDatabaseLogin    string `envconfig:"DATABASE_LOGIN" json:"database_login" default:"senya"`
	FlagDatabasePassword string `envconfig:"DATABASE_PASSWORD" json:"database_password" default:"1q2w3e4r5t"`
	FlagDatabaseName     string `envconfig:"DATABASE_NAME" json:"database_name" default:"MEDODS"`
	FlagDatabaseSSLMode  string `envconfig:"DATABASE_SSL_MODE" json:"database_ssl_mode" default:"disable"`
}
