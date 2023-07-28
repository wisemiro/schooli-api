package models

import "time"

type Shipping struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	User      User      `json:"user"`
	Geo       Geo       `json:"geo"`
}

type Geo struct {
	Latitude  float64 `json:"latitude" `
	Longitude float64 `json:"longitude"`
}
