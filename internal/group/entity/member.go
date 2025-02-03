package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type GroupMemberId string

type GroupMember struct {
	Id      GroupMemberId
	UserId  user_entity.UserId
	GroupId GroupId
}

func (id GroupMemberId) String() string {
	return string(id)
}
