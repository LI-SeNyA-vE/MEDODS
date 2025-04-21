package serverconfig

import "time"

func NewConfigServerMock() *Server {
	return &Server{
		FlagAddressAndPort: "localhost:8901",
		DatabaseCfg: DatabaseConfig{
			FlagDatabaseHost:     "localhost",
			FlagDatabasePort:     "5432",
			FlagDatabaseLogin:    "senya",
			FlagDatabasePassword: "1q2w3e4r5t",
			FlagDatabaseName:     "MEDODS",
			FlagDatabaseSSLMode:  "disable",
		},
		FlagAccessKey:       "f4pq3792h3dy4g82o63R84P265o3874wgfiy2p947gf7qo5hcnbvtbo8y2c9upnox3q9E3",
		FlagLifetimeAccess:  time.Minute * 15,
		FlagRefreshKey:      "x53416ucertiyvuybiunb5yp6no78iu65cr34exqyto839p28u320kjfnubviry3294bdf",
		FlagLifetimeRefresh: time.Hour * 24,
	}
}
