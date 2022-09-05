package gophermart

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/pkg/database"
	"go-gofermart-loyalty-system/internal/pkg/httpserver"
	"go-gofermart-loyalty-system/internal/router"
)

func Run(log *zap.Logger, cfg *config.Config) {
	log.Info("Connection to database...")
	db, err := database.New(cfg.DBUrl)

	if err != nil {
		log.Fatal("Can't connection to db", zap.Error(err))
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		log.Info("Close database connection")
		err := db.Close()
		if err != nil {
			log.Error("Error to close connection db", zap.Error(err))
		}
	}(db)

	log.Info("Staring the application...")

	apiMux := router.New(log)

	apiServer := httpserver.NewHttpServer(log, apiMux, cfg.Port)
	log.Info("Staring rest api server...")
	apiServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		log.Info("Received a signal.", zap.String("signal", x.String()))
	case err := <-apiServer.Notify():
		log.Error("Received an error from the start rest api server", zap.Error(err))
	}

	log.Info("Stopping server...")

	if err := apiServer.Stop(context.Background()); err != nil {
		log.Error("Got an error while stopping th rest api server", zap.Error(err))
	}

	log.Info("The gophermart is calling the last defers and will be stopped.")
}
