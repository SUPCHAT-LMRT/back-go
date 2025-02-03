package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_groups/chat_message/entity"
	group_chat_message_entity "github.com/supchat-lmrt/back-go/internal/group/chat_groups/entity"
)

type GroupChatMessageRepository interface {
	Create(ctx context.Context, message *entity.GroupChatMessage) error
	ListByGroupId(ctx context.Context, groupId group_chat_message_entity.ChatGroupId) ([]*entity.GroupChatMessage, error)
}
