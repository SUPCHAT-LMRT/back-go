package entity

import "time"

type InviteLink struct {
	Token     string
	FirstName string
	LastName  string
	Email     string
	ExpiresAt time.Time
}
