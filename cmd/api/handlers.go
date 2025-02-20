package main

import (
	"errors"
	"net/http"

	"anmol.gaud/internal/models"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addCacheEntryHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key string `json:"key"`
		Value string `json:"value"`
		TTL int64 `json:"ttl"`
	}
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestResponse(w, r, err);
		return
	}
	cacheEntry := &models.CacheEntry{
		Key: body.Key,
		Value: body.Value,
		TTL: body.TTL,
	}
	err = app.models.KeyValue.Insert(cacheEntry)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) getCacheEntryHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key string `json:"key"`
	}
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	cacheEntry, err := app.models.KeyValue.Get(body.Key)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.nilValueResponse(w, r)
		default:
			app.serverErrorResponse(w,r,err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"value": cacheEntry}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCacheEntryHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key string `json:"key"`
	}
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	err = app.models.KeyValue.Delete(body.Key)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			w.WriteHeader(http.StatusOK)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}