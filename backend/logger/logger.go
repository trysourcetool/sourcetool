package logger

import (
	"log"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"

	"github.com/trysourcetool/sourcetool/backend/config"
)

var Logger *zap.Logger

func Init() {
	var err error
	switch config.Config.Env {
	case config.EnvProd:
		Logger, err = zapdriver.NewProductionWithCore(zapdriver.WrapCore(
			zapdriver.ServiceName("server"),
		))
	case config.EnvStaging:
		Logger, err = zapdriver.NewDevelopmentWithCore(zapdriver.WrapCore(
			zapdriver.ServiceName("server"),
		))
	default:
		Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}

	defer Logger.Sync()
}
