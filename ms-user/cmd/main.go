package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rizface/go-ms-systemd/ms-user/cmd/rest"
	"github.com/rizface/go-ms-systemd/ms-user/database"
)

func runMigration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		os.Getenv("SYSTEMD_DB_NAME"),
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("✅ All migrations is up")
	}

	log.Println("✅ Success run migrations")

	return nil
}

func main() {
	ctx := context.Background()

	dbPool, err := database.NewPgConn(ctx)
	if err != nil {
		log.Fatalf("failed created db connection pool: %s", err)
	}

	db := stdlib.OpenDBFromPool(dbPool)

	err = runMigration(db)
	if err != nil {
		log.Fatalf("failed run migrations: %v", err)
	}

	log.Println("success establish connection with database")

	restServer := rest.NewServer(dbPool)
	if err := restServer.Start(); err != nil {
		log.Fatalf("failed start http server: %s", err)
	}
}
