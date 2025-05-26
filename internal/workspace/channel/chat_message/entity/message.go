package entity

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ChannelMessageId string

type ChannelMessage struct {
	Id        ChannelMessageId
	ChannelId channel_entity.ChannelId
	Content   string
	AuthorId  entity.UserId
	Reactions []*ChannelMessageReaction
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (id ChannelMessageId) String() string {
	return string(id)
}
