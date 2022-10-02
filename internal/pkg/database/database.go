package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func New(dburi string) (*sql.DB, error) {
	log.Print(dburi)
	db, err := sql.Open("pgx", dburi)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

type InitializerSchema interface {
	InitSchema(ctx context.Context) error
}

func InitSchemas(ctx context.Context, db *sql.DB, initialize ...InitializerSchema) error {
	for _, initializerSchema := range initialize {
		if err := initializerSchema.InitSchema(ctx); err != nil {
			return err
		}
	}

	return nil
}
