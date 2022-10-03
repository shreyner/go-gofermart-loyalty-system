package main

import (
	"fmt"
	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/gophermart"
	"go-gofermart-loyalty-system/internal/pkg/logger"
	"log"
)

func main() {
	cfg := &config.Config{}
	if err := cfg.Parse(); err != nil {
		log.Fatal("Can't parse env")
		fmt.Println(err)

		return
	}

	log, err := logger.InitLogger(cfg)
	if err != nil {
		log.Fatal("error initilizing logger")
		fmt.Println(err)

		return
	}

	defer log.Sync()

	gophermart.Run(log, cfg)
}
