package models

import "time"

type Notifications struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	DeviceID  int32     `json:"device_id"`
	UserID    int32     `json:"user_id"`
	Device    Devices   `json:"device"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Seen      bool      `json:"seen"`
}

type Devices struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int32     `json:"user_id"`
	Device    string    `json:"device"`
}
