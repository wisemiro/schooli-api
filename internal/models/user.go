package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"-"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	IsVerified   bool      `json:"is_verified"`
	Password     string
	PasswordHash string
}

func NewUser(email, password, phoneNumber string) (*User, error) {
	u := &User{
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
	}
	if err := u.Hash(); err != nil {
		return nil, fmt.Errorf("hashing password failed %w", err)
	}
	return u, nil
}

func (u *User) Hash() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	u.PasswordHash = string(bytes)
	return err
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type Payload struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	IssuedAt    time.Time `json:"issued_at"`
	ATExpiredAt time.Time `json:"expired_at"`
	RTExpiredAt time.Time `json:"rt_expired_at"`
}

// Valid checks if the token payload is valid or not.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ATExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}
