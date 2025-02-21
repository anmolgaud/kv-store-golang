package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anmol.gaud/internal/models"
	"anmol.gaud/internal/utils"
	"anmol.gaud/internal/yapper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dsn string
	}
}

type application struct {
	config config
	yapper *yapper.Logger
	models models.Models
	ticker time.Ticker
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "devlopment", "Environment (development|staging|production)")

	flag.Parse()

	yapper := yapper.New(os.Stdout, yapper.LevelInfo)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	db, err := openDB(cfg)
	if err != nil {
		yapper.PrintFatal(err, nil)
	}
	defer db.Close()
	err = applyMigrations(db)
	if err != nil {
		yapper.PrintFatal(err, nil)
	}
	yapper.PrintInfo("database connection pool established", nil)

	quit := make(chan os.Signal, 1)
	done := make(chan int)
	defer close(quit)
	defer close(done)

	app := &application{
		config: cfg,
		yapper: yapper,
		models: models.NewModel(db),
		ticker: *ticker,
	}
	
	go func(quit chan<- os.Signal) {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	}(quit)
	subscribers := utils.BroadCast(done, quit, 2)
	err = app.serve(quit, subscribers)
	if err != nil {
		done <- 1
		yapper.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err 
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err 
	}

	return db, nil 
}

func applyMigrations(db *sql.DB) error {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///Users/anmol/go-code/kv-store/migrations", "sqlite3", instance)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}