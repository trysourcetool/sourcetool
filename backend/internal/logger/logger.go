package logger

import (
	"log"

	"go.uber.org/zap"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

var Logger *zap.Logger

func Init() {
	var err error
	switch config.Config.Env {
	case config.EnvProd:
		Logger, err = zap.NewProduction()
	case config.EnvStaging:
		Logger, err = zap.NewProduction()
	default:
		Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}
}
