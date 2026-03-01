// Package postgres holds the pgxpool connection pool - internal/adapter/postgres/* build
// repositories on top of it.
package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	maxConns        = 10
	minConns        = 2
	maxConnLifetime = 5 * time.Minute
	maxConnIdleTime = 5 * time.Minute
)

type Config struct {
	User     string `envconfig:"POSTGRES_USER"     required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Port     string `envconfig:"POSTGRES_PORT"     required:"true"`
	Host     string `envconfig:"POSTGRES_HOST"     required:"true"`
	DBName   string `envconfig:"POSTGRES_DB_NAME"  required:"true"`
}

type Pool struct {
	*pgxpool.Pool
}

// New opens a connection pool and pings it, so a wrong password/host/etc. fails fast at startup
// instead of on the first query.
func New(ctx context.Context, c Config) (*Pool, error) {
	// Built as a URL rather than fmt.Sprintf'd keyword/value pairs so a password containing a
	// space, quote or other special character can't break the connection-string syntax - net/url
	// percent-encodes each component for us.
	dsn := (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   "/" + c.DBName,
	}).String()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	cfg.MaxConns = maxConns
	cfg.MinConns = minConns
	cfg.MaxConnLifetime = maxConnLifetime
	cfg.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()

	log.Info().Msg("postgres: closed")
}
