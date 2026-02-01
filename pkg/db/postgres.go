package db

import (
	"context"
	"fmt"
	"log"
	db_config "virtual_pet_game/config/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDSN(dbCfg *db_config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.SSLMode,
	)
}

func NewPool(ctx context.Context, dsn string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	return pool
}
