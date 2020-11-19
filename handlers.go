package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (app *App) Welcome(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"message": "Welcome to A-Reminder-App API"}`))

	return
}

func (app *App) CreateReminder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	body := &Reminder{}
	err = json.Unmarshal(b, body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	query := "INSERT INTO reminders (message, time, longitude, latitude) VALUES ($1, $2, $3, $4) RETURNING *"

	result := &Reminder{}
	err = app.db.
		QueryRow(query, body.Message, body.Time, body.Longitude, body.Latitude).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.CreatedAt, &result.UpdatedAt,
		)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	err = app.cache.Set(fmt.Sprintf("%s:%s", "reminders", result.ID), result)

	r := &APIResponse{Data: result}
	bb, err := json.Marshal(r)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(bb)

	return
}
