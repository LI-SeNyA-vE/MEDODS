package serverconfig

import (
	"github.com/kelseyhightower/envconfig"
)

func NewConfigServer() (*Server, error) {
	var cfg Server
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
