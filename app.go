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
func New(router *mux.Router, db *sqlx.DB, cache *CacheAPI) (app App, err error) {
	app.db = db
	app.cache = cache

	app.registerHandlers(router)

	return
}

func (app *App) welcome(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"message": "Welcome to A-Reminder-APP API"}`))

	return
}

func (app *App) registerHandlers(router *mux.Router) {
	rh := NewReminderHandler(app.db, app.cache)

	router.HandleFunc("/", app.welcome).Methods(http.MethodGet)
	router.HandleFunc("/api/v1", app.welcome).Methods(http.MethodGet)

	router.HandleFunc("/api/v1/reminders", rh.CreateReminder).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/reminders/{id}", rh.GetSingleReminder).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/reminders", rh.GetReminders).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/reminders/{id}", rh.UpdateReminder).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/reminders/{id}/status/complete", rh.MarkReminderAsComplete).Methods(http.MethodPatch)
	router.HandleFunc("/api/v1/reminders/{id}", rh.DeleteReminder).Methods(http.MethodDelete)
}
