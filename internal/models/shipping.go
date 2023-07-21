package models

import "time"

type Shipping struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
	Location    string    `json:"location"`
	Address     string    `json:"address"`
	Apartment   string    `json:"apartment"`
	PhoneNumber string    `json:"phone_number"`
	User        User      `json:"user"`
	Order       Order     `json:"orders"`
	Status      string    `json:"status"`
}

type Geo struct {
	Latitude  float64 `json:"latitude" `
	Longitude float64 `json:"longitude"`
}
