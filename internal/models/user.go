package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Password  string    `json:"-" db:"password"`
	Email     string    `json:"email" db:"email"`
	Country   string    `json:"country" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

//UsersList is the user list struct
type UsersList struct {
	Users      []*User `json:"users"`
	TotalCount uint64  `json:"total_count"`
}
