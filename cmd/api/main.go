package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"
	"os/signal"
	"syscall"
	"anmol.gaud/internal/models"
	"anmol.gaud/internal/yapper"
	"anmol.gaud/internal/utils"
	_ "github.com/mattn/go-sqlite3"
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
	db, err := openDB(cfg)
	if err != nil {
		yapper.PrintFatal(err, nil)
	}
	defer ticker.Stop()
	defer db.Close()
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
	db, err := sql.Open("sqlite3", "/Users/anmol/go-code/kv-store/data/store.db")
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