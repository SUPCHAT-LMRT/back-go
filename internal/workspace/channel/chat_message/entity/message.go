package entity

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ChannelMessageId string
type ChannelMessageAttachmentId string

type ChannelMessage struct {
	Id          ChannelMessageId
	ChannelId   channel_entity.ChannelId
	Content     string
	AuthorId    entity.UserId
	Reactions   []*ChannelMessageReaction
	Attachments []*ChannelMessageAttachment
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ChannelMessageAttachment struct {
	Id       ChannelMessageAttachmentId
	FileName string
}

func (id ChannelMessageId) String() string {
	return string(id)
}

func (c ChannelMessageAttachmentId) String() string {
	return string(c)
}
