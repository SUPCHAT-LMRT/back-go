package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type ChannelMessageReactionId string

type ChannelMessageReaction struct {
	Id        ChannelMessageReactionId
	MessageId ChannelMessageId
	UserId    user_entity.UserId
	Reaction  string
}

func (id ChannelMessageReactionId) String() string {
	return string(id)
}
