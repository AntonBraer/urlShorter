package postgres

import (
	"context"
	"fmt"
	"github.com/AntonBraer/urlShorter/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const (
	defaultConnectAttempts = 10
	defaultConnectTimeout  = time.Second
	maxPoolSize            = 1
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string, log *logger.Logger, ctx context.Context) (pg *Postgres, err error) {
	var poolConfig *pgxpool.Config
	if poolConfig, err = pgxpool.ParseConfig(url); err != nil {
		return nil, fmt.Errorf("pg - New - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = maxPoolSize

	attempts := 0
	for attempts < defaultConnectAttempts {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("pg - New - connection timeout exceeded")
		default:
			if pg.Pool, err = pgxpool.ConnectConfig(ctx, poolConfig); err == nil {
				return pg, nil
			}

			attempts++
			log.Info(fmt.Sprintf("Postgres is trying to connect, attempt = %d", attempts), attempts)
			time.Sleep(defaultConnectTimeout)
		}
	}

	return nil, fmt.Errorf("pg - New - pgxpool.ConnectConfig: %w", err)
}
