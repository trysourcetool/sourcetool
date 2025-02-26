package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}

func Sync() error {
	return Log.Sync()
}
