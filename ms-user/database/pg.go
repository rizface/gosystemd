package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgConn(ctx context.Context) (*pgxpool.Pool, error) {
	var (
		user     = "postgres"
		password = "postgres"
		port     = "5432"
		host     = "localhost"
		dbname   = "postgres"
	)

	if os.Getenv("SYSTEMD_DB_NAME") != "" {
		user = os.Getenv("SYSTEMD_DB_USER")
		password = os.Getenv("SYSTEMD_DB_PASSWORD")
		host = os.Getenv("SYSTEMD_DB_HOST")
		port = os.Getenv("SYSTEMD_DB_PORT")
		dbname = os.Getenv("SYSTEMD_DB_NAME")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dbconfig.MaxConnLifetime = 1 * time.Hour
	dbconfig.MaxConnIdleTime = 30 * time.Minute
	dbconfig.HealthCheckPeriod = 5 * time.Second
	dbconfig.MaxConns = 10
	dbconfig.MinConns = 5

	return pgxpool.NewWithConfig(ctx, dbconfig)
}
