package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// type APIHandlers interface {
// 	validator ValidationAPI
// }

type ReminderHandlers struct {
	validator *ReminderValidator
	service   *ServiceAPI
}

func NewReminderHandlers(db *sqlx.DB, cache *CacheAPI) *ReminderHandlers {
	validator, err := NewReminderValidator()
	if err != nil {
		log.Fatal(err)
	}
	service := NewServiceAPI(db, cache)

	rh := &ReminderHandlers{
		validator: validator,
		service:   service,
	}

	return rh
}

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

	rb, err := json.Marshal(result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	// save to cache
	err = app.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	r := &APIResponse{Data: result}
	bb, err := json.Marshal(r)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Write(bb)

	return
}

func (app *App) GetSingleReminder(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	result := &Reminder{}

	// check cache first
	s, err := app.cache.Get(fmt.Sprintf("reminders:%s", id))
	if err == nil {
		err := json.Unmarshal([]byte(s), result)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

			return
		}
		r := &APIResponse{Data: result}

		bb, err := json.Marshal(r)

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(bb))

		return
	}

	query := "SELECT * FROM reminders WHERE id=$1"
	err = app.db.QueryRow(query, id).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.CreatedAt, &result.UpdatedAt,
		)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": not found}`))

		return
	}

	// save to cache
	rb, err := json.Marshal(result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	// save to cache
	err = app.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

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
