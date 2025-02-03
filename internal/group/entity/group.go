package entity

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type GroupId string

type Group struct {
	Id          GroupId
	Name        string
	OwnerUserId user_entity.UserId
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (id GroupId) String() string {
	return string(id)
}
