package main

import "time"

type ReminderPayload struct {
	Message   string    `json:"message" validate:"required"`
	Time      time.Time `json:"time" validate:"required"`
	Latitude  float32   `json:"latitude" validate:"required,latitude"`
	Longitude float32   `json:"longitude" validate:"required,longitude"`
}

type Reminder struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
	Status    string    `json:"status"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type APIResponse struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}
