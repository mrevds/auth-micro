package entity

import "time"

type User struct {
	ID        string
	Username  string
	Name      string
	Email     string
	Age       int32
	Bio       string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
