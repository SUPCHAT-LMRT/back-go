package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type ChatGroupMemberId string

type ChatGroupMember struct {
	Id     ChatGroupId
	UserId user_entity.UserId
}
