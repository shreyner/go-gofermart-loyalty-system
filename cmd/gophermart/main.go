package main

import (
	"fmt"
	"os"

	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/gophermart"
	"go-gofermart-loyalty-system/internal/pkg/logger"
)

func main() {
	cfg := &config.Config{}
	if err := cfg.Parse(); err != nil {
		fmt.Errorf("Can't parse env: %w", err)
		os.Exit(1)
	}

	log, err := logger.InitLogger(cfg)
	if err != nil {
		fmt.Errorf("error initilizing logger: %w", err)
		os.Exit(1)
	}

	defer log.Sync()

	gophermart.Run(log, cfg)

}
