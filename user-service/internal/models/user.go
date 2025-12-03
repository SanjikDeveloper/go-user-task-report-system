package models

import "time"

type User struct {
	Id         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password"`
	Role       string    `json:"role" db:"role"`
	DateCreate time.Time `json:"date_create" db:"date_create"`
}
