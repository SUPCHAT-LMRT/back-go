package entity

import (
	"time"
)

type UserId string

type User struct {
	Id                   UserId
	FirstName            string
	LastName             string
	Email                string
	Password             string
	NotificationsEnabled bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (id UserId) String() string {
	return string(id)
}

func (id UserId) IsAfter(other UserId) bool {
	return id > other
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
