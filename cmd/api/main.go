package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/jsonlog"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

type envelope map[string]any
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "Database connection URL")

	flag.Parse()
	db, err := openDB(cfg)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	logger.PrintInfo("database connection Established", nil)

	app := &application{
		config: cfg,
		logger: logger,
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	app.logger.PrintInfo(fmt.Sprintf("App is running on http://localhost:%d", app.config.port), nil)
	err = srv.ListenAndServe()
	app.logger.PrintFatal(err, nil)
}
func (app *application) demo(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World")
	if err != nil {
		return
	}
}

// openDB establishes a connection to a PostgreSQL database using the provided configuration.
//
// Parameters:
//   - cfg: A `config` struct that contains the configuration needed to connect to the database.
//     The `cfg.db.dsn` field should contain the Data Source Name (DSN) for the PostgreSQL database.
//
// Returns:
// - (*sql.DB): A pointer to the `sql.DB` object representing the connection to the database, or nil if an error occurs.
// - error: An error if the connection could not be established or if the `PingContext` check fails, or nil if successful.
//
// This function opens a connection to a PostgreSQL database, verifies the connection by pinging the database,
// and returns the `sql.DB` object if successful. If there is an error during connection or pinging,
// it returns the error to allow the caller to handle it appropriately.
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
