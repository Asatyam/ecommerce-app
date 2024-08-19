package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/jsonlog"
	"github.com/Asatyam/ecommerce-app/internal/mailer"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"sync"
	"time"
)

type envelope map[string]any
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "Database connection URL")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP server hostname")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP Port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "26aac60d0744ad", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "705980f9080f4f", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Satyam Agrawal <agrasatyam1282@gmail.com>", "SMTP sender")
	
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
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
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
