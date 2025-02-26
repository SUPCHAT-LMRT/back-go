package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ChannelMessageRepository interface {
	Create(ctx context.Context, message *entity.ChannelMessage) error
	ListByChannelId(ctx context.Context, channelId channel_entity.ChannelId) ([]*entity.ChannelMessage, error)
	AddReaction(ctx context.Context, reaction entity.ChannelMessageReaction) error
}
