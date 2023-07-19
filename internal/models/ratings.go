package models

import "time"

type ProductRatings struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	Stars     int       `json:"stars"`
	Feeedback string    `json:"feeedback"`
}
