package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	Incomplete = "incomplete"
	Complete   = "complete"
)

type ReminderHandler struct {
	validator *ReminderValidator
	service   *ServiceAPI
}

func NewReminderHandler(db *sqlx.DB, cache *CacheAPI) *ReminderHandler {
	validator, err := NewReminderValidator()
	if err != nil {
		log.Fatal(err)
	}
	service := NewServiceAPI(db, cache)

	rh := &ReminderHandler{
		validator: validator,
		service:   service,
	}

	return rh
}

func (rh *ReminderHandler) CreateReminder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	body := &ReminderPayload{}
	err = json.Unmarshal(b, body)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	vErr := rh.validator.Validate(body)
	if vErr != nil {
		b, _ := json.Marshal(vErr)
		res.WriteHeader(http.StatusBadRequest)
		res.Write(b)

		return
	}

	r := Reminder{
		Message:   body.Message,
		Time:      body.Time,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	}
	result, err := rh.service.CreateReminder(r)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	rs := &APIResponse{Data: result}
	bb, err := json.Marshal(rs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Write(bb)

	return
}

func (rh *ReminderHandler) GetSingleReminder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id := vars["id"]

	result, err := rh.service.GetReminderByID(id)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": not found}`))

		return
	}

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

func (rh *ReminderHandler) GetReminders(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	statusQuery := strings.ToLower(req.URL.Query().Get("status"))
	limitQuery := strings.ToLower(req.URL.Query().Get("limit"))
	offsetQuery := strings.ToLower(req.URL.Query().Get("offset"))

	if limitQuery == "" {
		limitQuery = "10"
	}
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	var err error
	if statusQuery != "" && statusQuery != Incomplete && statusQuery != Complete {
		err = errors.New("status can either be incomplete or complete ")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{"error": status can either be incomplete or complete }`))

		return
	}

	limit := 10
	offset := 0

	if limit, err = strconv.Atoi(limitQuery); err != nil {
		err = errors.New("limit and offset must be positive integers")

		return
	}
	if offset, err = strconv.Atoi(offsetQuery); err != nil {
		err = errors.New("limit and offset must be positive integers")

		return
	}

	var result []*Reminder
	var count int
	if statusQuery != "" {
		if count, err = rh.service.GetRemindersCount("WHERE status=$1", statusQuery); err == nil {
			query := "SELECT * FROM reminders WHERE status=$1 LIMIT $2 OFFSET $3;"
			result, err = rh.service.GetReminders(query, statusQuery, limit, offset)
		}
	} else if count, err = rh.service.GetRemindersCount(""); err == nil {
		query := "SELECT * FROM reminders LIMIT $1 OFFSET $2;"
		result, err = rh.service.GetReminders(query, limit, offset)
	}

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"error": Something went wrong}`))

		return
	}

	pagination := Pagination{
		Limit:  limit,
		Offset: offset,
		Total:  count,
	}

	r := &APIResponse{Data: result, Pagination: pagination}
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

func (rh *ReminderHandler) UpdateReminder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id := vars["id"]

	// read payload
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	body := &ReminderPayload{}
	err = json.Unmarshal(b, body)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	// validate payload
	vErr := rh.validator.Validate(body)
	if vErr != nil {
		b, _ := json.Marshal(vErr)
		res.WriteHeader(http.StatusBadRequest)
		res.Write(b)

		return
	}

	// update data
	r := Reminder{
		Message:   body.Message,
		Time:      body.Time,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	}
	result, err := rh.service.UpdateReminder(id, r)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	rs := &APIResponse{Data: result}
	bb, err := json.Marshal(rs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(bb)

	return
}

func (rh *ReminderHandler) MarkReminderAsComplete(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id := vars["id"]

	result, err := rh.service.UpdateReminderStatus(id, Complete)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	rs := &APIResponse{Data: result}
	bb, err := json.Marshal(rs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(bb)

	return
}

func (rh *ReminderHandler) DeleteReminder(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id := vars["id"]

	err := rh.service.DeleteReminder(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))

		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("{ success: true }"))

	return
}
