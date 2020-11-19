package main

import "time"

type Reminder struct {
	ID        string    `json:"id,omit_empty"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at,omit_empty"`
}

type APIResponse struct {
	Data interface{} `json:"data"`
}
