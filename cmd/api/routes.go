package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/key", app.getCacheEntryHandler)
	router.HandlerFunc(http.MethodPost, "/v1/key",  app.addCacheEntryHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/key",  app.deleteCacheEntryHandler)


	return app.recoverPanic(router);
}