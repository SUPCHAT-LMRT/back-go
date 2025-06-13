package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type MessageReactionId string

type MessageReaction struct {
	Id       MessageReactionId
	UserIds  []user_entity.UserId
	Reaction string
}

func (id MessageReactionId) String() string {
	return string(id)
}
