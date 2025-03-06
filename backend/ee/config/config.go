package config

import (
	"log"

	"github.com/caarlos0/env/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
)

var Config *ConfigEE

const (
	EnvLocal = "local"
	EnvStg   = "stg"
	EnvProd  = "prod"
)

type ConfigEE struct {
	config.ConfigCE
}

func Init() {
	cfg := new(ConfigEE)
	envOpts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, envOpts); err != nil {
		log.Fatal("[INIT] config: ", err)
	}

	Config = cfg
}
