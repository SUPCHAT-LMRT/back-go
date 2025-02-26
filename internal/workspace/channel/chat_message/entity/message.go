package entity

import (
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"time"
)

type ChannelMessageId string

type ChannelMessage struct {
	Id        ChannelMessageId
	ChannelId channel_entity.ChannelId
	Content   string
	AuthorId  entity.UserId
	CreatedAt time.Time
	Reactions []*ChannelMessageReaction
}

func (id ChannelMessageId) String() string {
	return string(id)
}
