package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type GroupChatMessageRepository interface {
	Create(ctx context.Context, message *entity.GroupChatMessage) error
	ListByGroupId(ctx context.Context, groupId group_entity.GroupId) ([]*entity.GroupChatMessage, error)
}
