package repository

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ChannelMessageRepository interface {
	Create(ctx context.Context, message *entity.ChannelMessage) error
	ListByChannelId(ctx context.Context, channelId channel_entity.ChannelId) ([]*entity.ChannelMessage, error)
	ToggleReaction(ctx context.Context, messageId entity.ChannelMessageId, userId user_entity.UserId, reaction string) (added bool, err error)
}
