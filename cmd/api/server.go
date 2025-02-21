package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"errors"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <- quit
		app.yapper.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.yapper.PrintInfo("starting server", map[string]string{
		"address": srv.Addr,
		"environment": app.config.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err;
	}

	err = <- shutdownError
	if err != nil {
		return err
	}

	app.yapper.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}