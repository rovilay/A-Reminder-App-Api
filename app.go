package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// App ...
type App struct {
	db    *sqlx.DB
	cache *CacheAPI
}

// New returns new app instance
func New(router *mux.Router, db *sqlx.DB, cacheApi *CacheAPI) (app App, err error) {
	app.db = db
	app.cache = cacheApi

	app.registerHandlers(router)

	return
}

func (app *App) registerHandlers(router *mux.Router) {
	router.HandleFunc("/", app.Welcome).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/reminders", app.CreateReminder).Methods(http.MethodPost)
}
