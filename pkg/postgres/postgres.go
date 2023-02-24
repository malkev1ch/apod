package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	maxConns          = 60
	maxConnLifetime   = 120
	maxConnIdleTime   = 20
	connectionTimeout = 10 * time.Second
)

// NewPool returns new Postgresql pool object.
func NewPool(connStr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	// init postgres database pool
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgesql connection config: %w", err)
	}

	cfg.MaxConns = maxConns
	cfg.MaxConnLifetime = maxConnLifetime
	cfg.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to postgresql: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to get response from postgresql: %w", err)
	}

	return pool, nil
}
