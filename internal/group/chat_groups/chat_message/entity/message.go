package entity

import (
	chat_group_entity "github.com/supchat-lmrt/back-go/internal/group/chat_groups/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type GroupChatMessageId string

type GroupChatMessage struct {
	Id        GroupChatMessageId
	GroupId   chat_group_entity.ChatGroupId
	Content   string
	AuthorId  entity.UserId
	CreatedAt time.Time
}
