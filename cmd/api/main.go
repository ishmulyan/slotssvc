package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	log := log.New(os.Stdout, "slotssrv: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Configuration
	cfg := struct {
		Host            string        `default:"0.0.0.0:3000" envconfig:"HOST"`
		ReadTimeout     time.Duration `default:"5s" envconfig:"READ_TIMEOUT"`
		WriteTimeout    time.Duration `default:"5s" envconfig:"WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
	}{}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Parsing config: %v", err)
	}

	srv := http.Server{
		Addr:         cfg.Host,
		Handler:      handler(log),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("Server listening: %s", cfg.Host)
		serverErrors <- srv.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-osSignals:
		log.Println("Start shutdown...")

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", cfg.ShutdownTimeout, err)
			if err := srv.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}
}
