package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"errors"
)

func (app *application) serve(inputs []chan os.Signal) error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	
	go func(subscribe <-chan os.Signal) {
		s := <- subscribe
		app.yapper.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}(inputs[0])

	go func(subscribe <-chan os.Signal){
		for {
			select {
			case <-app.ticker.C:
				err := app.CleanUp()
				if err != nil {
					app.yapper.PrintError(err, map[string]string{"message": "Failed to cleanup"})
					shutdownError <- err
					return
				}
			case <-subscribe:
				app.yapper.PrintInfo("exiting cleanup goroutine", nil)
				return
			}
		}
	}(inputs[1])

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