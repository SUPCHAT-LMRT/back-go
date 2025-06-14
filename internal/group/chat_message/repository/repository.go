package repository

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ChatMessageRepository interface {
	Create(ctx context.Context, chatMessage *entity.GroupChatMessage) error
	Get(ctx context.Context, chatMessageId entity.GroupChatMessageId) (*entity.GroupChatMessage, error)
	GetLastMessage(ctx context.Context, groupId group_entity.GroupId) (*entity.GroupChatMessage, error)
	IsFirstMessage(ctx context.Context, groupId group_entity.GroupId) (bool, error)
	ListMessages(
		ctx context.Context,
		groupId group_entity.GroupId,
		params ListMessagesQueryParams,
	) ([]*entity.GroupChatMessage, error)
	ToggleReaction(
		ctx context.Context,
		messageId entity.GroupChatMessageId,
		userId user_entity.UserId,
		reaction string,
	) (added bool, err error)
	DeleteMessage(ctx context.Context, messageId entity.GroupChatMessageId) error
	UpdateMessage(
		ctx context.Context,
		message *entity.GroupChatMessage,
	) error
}

type ListMessagesQueryParams struct {
	Limit           int
	Before          time.Time
	After           time.Time
	AroundMessageId entity.GroupChatMessageId
}
