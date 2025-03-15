package user

import (
	"time"
)

type SearchUser struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
