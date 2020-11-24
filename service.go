package main

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ServiceAPI struct {
	db    *sqlx.DB
	cache *CacheAPI
}

func NewServiceAPI(db *sqlx.DB, cache *CacheAPI) *ServiceAPI {
	return &ServiceAPI{db: db, cache: cache}
}

func (s *ServiceAPI) CreateReminder(data Reminder) (result *Reminder, err error) {
	query := "INSERT INTO reminders (message, time, longitude, latitude) VALUES ($1, $2, $3, $4) RETURNING *"
	err = s.db.
		QueryRow(query, data.Message, data.Time, data.Longitude, data.Latitude).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.CreatedAt, &result.UpdatedAt,
		)

	if err != nil {
		return
	}

	rb, err := json.Marshal(result)
	if err != nil {
		return
	}

	// save to cache
	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}

func (s *ServiceAPI) GetReminderByID(id string) (result *Reminder, err error) {
	// check cache first
	res, err := s.cache.Get(fmt.Sprintf("reminders:%s", id))
	if err == nil {
		err = json.Unmarshal([]byte(res), result)
		if err != nil {
			return
		}

		return
	}

	query := "SELECT * FROM reminders WHERE id=$1"
	err = s.db.QueryRow(query, id).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.CreatedAt, &result.UpdatedAt,
		)
	if err != nil {
		return
	}

	// save to cache
	rb, err := json.Marshal(result)
	if err != nil {
		return
	}

	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}
