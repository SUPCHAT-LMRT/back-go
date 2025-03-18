package user

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type SearchUser struct {
	Id        user_entity.UserId
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
