package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/handling"
	"github.com/heimdal-rw/chmgt/models"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	handler, err := handling.NewHandler(config)
	if err != nil {
		log.Fatal(err)
	}
	handler.Datasource, err = models.NewDatasource(
		config.DatabaseConnection(),
		config.Database.Name,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Datasource.Close()

	srv := &http.Server{
		Handler:      handler.Router,
		Addr:         config.ListenAddr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Recevied shutdown request")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		handler.Datasource.Close()
		srv.Shutdown(ctx)
		os.Exit(0)
	}()

	log.Printf("Listening on %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
