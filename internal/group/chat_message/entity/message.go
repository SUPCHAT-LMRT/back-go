package entity

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type GroupChatMessageId string

type GroupChatMessage struct {
	Id        GroupChatMessageId
	GroupId   group_entity.GroupId
	Content   string
	AuthorId  entity.UserId
	CreatedAt time.Time
}

func (id GroupChatMessageId) String() string {
	return string(id)
}
