package user

import (
	"time"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type SearchUser struct {
	Id        user_entity.UserId
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
