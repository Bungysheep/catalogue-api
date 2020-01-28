package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
	"github.com/bungysheep/catalogue-api/pkg/protocols/rest"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := startUp(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start http server, error: %v", err)
	}
}

func startUp() error {
	restServer := rest.NewRestServer()

	if err := database.CreateDbConnection(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			ctx := context.TODO()

			log.Printf("Closing database connection...\n")
			database.DbConnection.Close()

			log.Printf("Stoping http server...\n")
			restServer.Shutdown(ctx)
		}
	}()

	log.Printf("Starting http server...\n")
	return restServer.RunServer()
}
