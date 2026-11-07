package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// ---------- Setup: version and types ----------
// This section declares what the program needs before any logic runs. The
// version constant is a fixed string, usually reported by a healthcheck
// endpoint. The config struct holds settings (port and environment) that are
// filled from command-line flags. The application struct bundles the shared
// dependencies (config and logger) so handlers can be written as methods on it
// and reach them without globals.

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// ---------- Reading configuration ----------
	// Here we create an empty config and fill it from command-line flags.
	// IntVar and StringVar register -port (default 4000) and -env (default
	// "development") and bind them to the struct fields. flag.Parse() then
	// reads the actual command line; without it, the flags above would never
	// take effect.
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// ---------- Building the app's dependencies ----------
	// We create a structured logger that writes readable text to stdout (nil
	// means default options, Info level and above). Then we build the
	// application struct, injecting the config and logger. It is a pointer so
	// every method shares the same single instance.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	// ---------- Configuring the HTTP server ----------
	// This block defines how the server behaves. Addr sets the listening
	// address (e.g. ":4000" for all interfaces). Handler is the router from
	// app.routes() (defined in another file), which decides what serves each
	// request. IdleTimeout, ReadTimeout and WriteTimeout cap how long
	// keep-alive connections, request reads and response writes may take, which
	// protects against slow or stuck clients. ErrorLog routes the server's
	// internal errors through slog at Error level so all logs share one format.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// ---------- Starting the server and handling failure ----------
	// First we log a startup message with the address and environment. Then
	// ListenAndServe blocks while the server runs; it only returns when
	// something goes wrong or the server stops, and the error is never nil. We
	// log that error and exit with status 1, signalling failure to the OS or
	// process manager.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
