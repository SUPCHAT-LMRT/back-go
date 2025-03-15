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
	UpdatedAt time.Time
}

type UserStatus string

var (
	UserStatusOnline  UserStatus = "online"
	UserStatusIdle    UserStatus = "idle"
	UserStatusOffline UserStatus = "offline"
)

func (id UserId) String() string {
	return string(id)
}

func (id UserId) IsAfter(other UserId) bool {
	return id > other
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
