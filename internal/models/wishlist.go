package models

import "time"

type Wishlist struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	Product   Product   `json:"product"`
	User      User      `json:"user"`
}
