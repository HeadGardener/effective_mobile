package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeadGardener/effective_mobile/internal/client"
	"github.com/HeadGardener/effective_mobile/internal/config"
	"github.com/HeadGardener/effective_mobile/internal/handlers"
	"github.com/HeadGardener/effective_mobile/internal/server"
	"github.com/HeadGardener/effective_mobile/internal/services"
	"github.com/HeadGardener/effective_mobile/internal/storage"
)

const shutdownTimeout = 5 * time.Second

var confPath = flag.String("conf-path", "./config/.env", "path to config env")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.Init(*confPath)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while initializing config: %s", err.Error())
	}

	db, err := storage.NewDB(conf.DBConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %s", err.Error())
	}

	var (
		personStorage = storage.NewPersonStorage(db)
	)

	var (
		httpClient = client.NewClient(conf.HTTPClientConfig)
	)

	var (
		personService = services.NewPersonService(personStorage, httpClient)
	)

	handler := handlers.NewHandler(personService)

	srv := &server.Server{}
	go func() {
		if err = srv.Run(conf.ServerConfig, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %s", err.Error())
		}
	}()
	log.Println("[INFO] server start working")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("[INFO] server forced to shutdown: %e", err)
	}

	if err = db.Close(); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %e", err)
	}

	log.Println("[INFO] server exiting")
}
