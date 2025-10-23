package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	kvstore "github.com/shardul-d/kv-store"
	"github.com/tidwall/redcon"
)

var (
	// Version of the build. This is injected at build-time.
	buildString = "unknown"
)

type App struct {
	lo      *slog.Logger
	kvstore *kvstore.KVStore
}

func main() {
	// Initialise and load the config.
	ko, err := initConfig()
	if err != nil {
		fmt.Printf("error loading config: %v", err)
		os.Exit(-1)
	}

	app := &App{
		lo: initLogger(ko),
	}
	app.lo.Info("booting kvstoredb server", "version", buildString)

	// Set config options for kvstore.
	cfg := []kvstore.Config{kvstore.WithDir(ko.MustString("app.dir")), kvstore.WithAutoSync()}
	if ko.Bool("app.read_only") {
		cfg = append(cfg, kvstore.WithReadOnly())
	}
	if ko.Bool("app.debug") {
		cfg = append(cfg, kvstore.WithDebug())
	}

	// Initialise kvstore.
	kvstore, err := kvstore.Init(cfg...)
	if err != nil {
		app.lo.Error("error opening kvstore db", "error", err)
	}
	app.kvstore = kvstore

	// Initialise server.
	mux := redcon.NewServeMux()
	mux.HandleFunc("ping", app.ping)
	mux.HandleFunc("quit", app.quit)
	mux.HandleFunc("set", app.set)
	mux.HandleFunc("get", app.get)
	mux.HandleFunc("del", app.delete)

	// Create a channel to listen for cancellation signals.
	// Create a new context which is cancelled when `SIGINT`/`SIGTERM` is received.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	srvr := redcon.NewServer(ko.MustString("server.address"),
		mux.ServeRESP,
		func(conn redcon.Conn) bool {
			// use this function to accept or deny the connection.
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
		},
	)

	// Sart the server in a goroutine.
	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			app.lo.Error("failed to listen and serve", "error", err)
		}
	}()

	// Listen on the close channel indefinitely until a
	// `SIGINT` or `SIGTERM` is received.
	<-ctx.Done()

	// Cancel the context to gracefully shutdown and perform
	// any cleanup tasks.
	cancel()
	app.kvstore.Shutdown()
	srvr.Close()
}
