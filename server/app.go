package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type App struct{}

func New(router *mux.Router) (app *App, err error) {
	app.registerHandlers(router)

	return
}

func (app *App) registerHandlers(router *mux.Router) {
	router.HandleFunc("/", app.Welcome).Methods(http.MethodGet)
}
