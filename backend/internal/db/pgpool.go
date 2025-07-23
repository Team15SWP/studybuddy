package db

import (
	"context"
	"fmt"

	"study_buddy/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgPool(ctx context.Context, pgConfig *config.PGConfig) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgConfig.DBUser, pgConfig.DBPass, pgConfig.DBHost, pgConfig.DBPort, pgConfig.DBName))
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	err = pgPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pgPool.Ping: %w", err)
	}

	return pgPool, nil
}
