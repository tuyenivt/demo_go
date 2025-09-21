package models

import (
	"time"
)

type User struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	PasswordHashed string     `json:"-"`
	LastLogin      *time.Time `json:"last_login"`
	TimeCreated    time.Time  `json:"time_created"`
	TimeConfirmed  *time.Time `json:"time_confirmed"`
	TimeDeleted    *time.Time `json:"time_deleted"`
	Favorites      []Movie    `json:"favorites"`
	Watchlist      []Movie    `json:"watchlist"`
}
