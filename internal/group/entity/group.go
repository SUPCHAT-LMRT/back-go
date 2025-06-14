package entity

import (
	"time"
)

type GroupId string

type Group struct {
	Id            GroupId
	Name          string
	OwnerMemberId GroupMemberId
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (id GroupId) String() string {
	return string(id)
}
