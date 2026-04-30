package model

import "time"

type Hospital struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	APIURL    string    `json:"api_url" db:"api_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
