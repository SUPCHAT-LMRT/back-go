package entity

import (
	"time"
)

type UserId string

type User struct {
	Id         UserId
	FirstName  string
	LastName   string
	Email      string
	Pseudo     string
	Password   string
	IsVerified bool // IsVerified is true if the user has validated his email address
	BirthDate  time.Time
	CreatedAt  time.Time
}

func (id UserId) String() string {
	return string(id)
}
