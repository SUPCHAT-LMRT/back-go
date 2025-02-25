package entity

import (
	"time"
)

type UserId string

type User struct {
	Id        UserId
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (id UserId) String() string {
	return string(id)
}
