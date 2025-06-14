package repository

import (
	"context"
	"time"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ChannelMessageRepository interface {
	Create(ctx context.Context, message *entity.ChannelMessage) error
	Get(ctx context.Context, id entity.ChannelMessageId) (*entity.ChannelMessage, error)
	ListByChannelId(
		ctx context.Context,
		channelId channel_entity.ChannelId,
		params ListByChannelIdQueryParams,
	) ([]*entity.ChannelMessage, error)
	ToggleReaction(
		ctx context.Context,
		messageId entity.ChannelMessageId,
		userId user_entity.UserId,
		reaction string,
	) (added bool, err error)
	CountByWorkspace(ctx context.Context, id workspace_entity.WorkspaceId) (uint, error)
	ListAllMessagesByUser(ctx context.Context, userId user_entity.UserId) ([]*entity.ChannelMessage, error)
	DeleteMessage(ctx context.Context, channelMessageId entity.ChannelMessageId) error
	UpdateMessage(ctx context.Context, msg *entity.ChannelMessage) error
}

type ListByChannelIdQueryParams struct {
	Limit           int
	Before          time.Time
	After           time.Time
	AroundMessageId entity.ChannelMessageId
}
