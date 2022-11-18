package models

import "time"

type RecordClient struct {
	Protocol string
	Hostname string
}

type TodoItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
