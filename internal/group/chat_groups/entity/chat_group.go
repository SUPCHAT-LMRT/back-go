package entity

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ChatGroupId string

type ChatGroup struct {
	Id          ChatGroupId
	GroupId     group_entity.GroupId
	OwnerUserId user_entity.UserId
	Members     []*ChatGroupMember
}
