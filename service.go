package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	query := "INSERT INTO reminders (message, time, longitude, latitude) VALUES ($1, $2, $3, $4) RETURNING *;"
	result = &Reminder{}
	err = s.db.
		QueryRow(query, data.Message, data.Time, data.Longitude, data.Latitude).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.Status, &result.CreatedAt, &result.UpdatedAt,
		)

	if err != nil {
		log.Println("ERROR > CreateReminder: ", err)
		return
	}

	rb, err := json.Marshal(result)
	if err != nil {
		log.Println("ERROR > CreateReminder: ", err)
		return
	}

	// save to cache
	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}

func (s *ServiceAPI) GetReminderByID(id string) (result *Reminder, err error) {
	result = &Reminder{}

	// check cache first
	res, err := s.cache.Get(fmt.Sprintf("reminders:%s", id))
	if err == nil {
		err = json.Unmarshal([]byte(res), result)
		if err != nil {
			log.Println("ERROR > GetReminderByID: ", err)
			return
		}

		return
	}

	query := "SELECT * FROM reminders WHERE id=$1;"
	err = s.db.QueryRow(query, id).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.Status, &result.CreatedAt, &result.UpdatedAt,
		)
	if err != nil {
		log.Println("ERROR > GetReminderByID: ", err)
		return
	}

	// save to cache
	rb, err := json.Marshal(result)
	if err != nil {
		log.Println("ERROR > GetReminderByID: ", err)
		return
	}

	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}

func (s *ServiceAPI) GetRemindersCount(whereClause string, args ...interface{}) (count int, err error) {
	query := `
		SELECT
			CAST(COUNT(id) AS INTEGER)
		FROM
			reminders
		`
	if len(args) > 0 {
		query = fmt.Sprintf("%s %s;", query, whereClause)
		err = s.db.QueryRow(query, args...).Scan(&count)
	} else {
		err = s.db.QueryRow(fmt.Sprintf("%s;", query)).Scan(&count)
	}

	if err != nil {
		log.Println("ERROR > GetRemindersCount: ", err)
		return
	}

	return
}

func (s *ServiceAPI) GetReminders(query string, args ...interface{}) (result []*Reminder, err error) {
	var rows *sql.Rows

	if len(args) > 0 {
		rows, err = s.db.Query(query, args...)
	} else {
		rows, err = s.db.Query(query)
	}

	if err != nil {
		log.Println("ERROR > GetReminders: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		res := &Reminder{}
		err := rows.Scan(
			&res.ID, &res.Message, &res.Time,
			&res.Longitude, &res.Latitude, &res.Status, &res.CreatedAt, &res.UpdatedAt,
		)
		if err != nil {
			log.Println("ERROR > GetReminders: ", err)
			return result, err
		}

		result = append(result, res)
	}
	// Check for errors from iterating over rows.
	err = rows.Err()

	return
}

func (s *ServiceAPI) UpdateReminder(id string, data Reminder) (result *Reminder, err error) {
	query := "UPDATE reminders SET message=$1, time=$2, longitude=$3, latitude=$4, updated_at=$5 WHERE id=$6 AND status!=$7 RETURNING *;"
	result = &Reminder{}

	err = s.db.
		QueryRow(query, data.Message, data.Time, data.Longitude, data.Latitude, time.Now(), id, Complete).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.Status, &result.CreatedAt, &result.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("no record found or updated")
		}
		log.Println("ERROR > UpdateReminder: ", err)
		return
	}

	rb, err := json.Marshal(result)
	if err != nil {
		log.Println("ERROR > UpdateReminder: ", err)
		return
	}

	// save to cache
	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}

func (s *ServiceAPI) UpdateReminderStatus(id string, status string) (result *Reminder, err error) {
	query := "UPDATE reminders SET status=$1, updated_at=$2 WHERE id=$3 AND status!=$4 RETURNING *;"
	result = &Reminder{}

	err = s.db.
		QueryRow(query, status, time.Now(), id, status).
		Scan(
			&result.ID, &result.Message, &result.Time,
			&result.Longitude, &result.Latitude, &result.Status, &result.CreatedAt, &result.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("no record found or updated")
		}

		log.Println("ERROR > UpdateReminderStatus: ", err)
		return
	}

	rb, err := json.Marshal(result)
	if err != nil {
		log.Println("ERROR > UpdateReminderStatus: ", err)
		return
	}

	// update cache
	err = s.cache.Set(fmt.Sprintf("reminders:%s", result.ID), rb)

	return
}

func (s *ServiceAPI) DeleteReminder(id string) (err error) {
	// remove from cache
	s.cache.Del(fmt.Sprintf("reminders:%s", id))

	// remove from db
	query := "DELETE FROM reminders WHERE id=$1;"

	_, err = s.db.Exec(query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("no record found or deleted")
		}

		log.Println("ERROR > UpdateReminderStatus: ", err)
	}

	return
}
