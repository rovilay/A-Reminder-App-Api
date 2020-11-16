package main

import "net/http"

func (app *App) Welcome(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"message": "Welcome to A-Reminder-App API"}`))

	return
}
