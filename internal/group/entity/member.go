package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type GroupMemberId string

type GroupMember struct {
	Id     GroupMemberId
	UserId user_entity.UserId
}

func (id GroupMemberId) String() string {
	return string(id)
}
