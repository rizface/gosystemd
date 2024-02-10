package main

import (
	"context"
	"log"

	"github.com/rizface/go-ms-systemd/ms-user/cmd/rest"
	"github.com/rizface/go-ms-systemd/ms-user/database"
)

func main() {
	ctx := context.Background()

	dbPool, err := database.NewPgConn(ctx)
	if err != nil {
		log.Fatalf("failed created db connection pool: %s", err)
	}

	restServer := rest.NewServer(dbPool)
	if err := restServer.Start(); err != nil {
		log.Fatalf("failed start http server: %s", err)
	}
}
