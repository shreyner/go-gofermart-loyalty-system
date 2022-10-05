package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"go.uber.org/zap"
	"log"

	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/gophermart"
	"go-gofermart-loyalty-system/internal/pkg/logger"
)

func migration(ctx context.Context, log *zap.Logger, databaseURI string) error {
	conn, err := pgx.Connect(ctx, databaseURI)

	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)

	if err != nil {
		return err
	}

	migration, err := migrate.NewMigrator(ctx, conn, "public.schema_version")

	if err != nil {
		return err
	}

	err = migration.LoadMigrations("migrate")

	if err != nil {
		return err
	}

	if len(migration.Migrations) == 0 {
		return nil
	}

	err = migration.Migrate(ctx)

	if err != nil {
		return err
	}

	return nil
}

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

	log.Info("Start migrations ...")
	if err = migration(context.Background(), log, cfg.DBURI); err != nil {
		log.Error("error migration", zap.Error(err))
		fmt.Println(err)

		return
	}

	log.Info("Finish migrations")

	gophermart.Run(log, cfg)
}
