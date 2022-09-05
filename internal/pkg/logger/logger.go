package logger

import (
	"fmt"

	"go-gofermart-loyalty-system/internal/config"

	"go.uber.org/zap"
)

func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	cfgLog := zap.NewProductionConfig()

	if cfg.IsDev {
		cfgLog = zap.NewDevelopmentConfig()
	}

	cfgLog.DisableStacktrace = true

	logger, err := cfgLog.Build()

	if err != nil {
		return nil, fmt.Errorf("logger build: %w", err)
	}

	return logger, nil
}
