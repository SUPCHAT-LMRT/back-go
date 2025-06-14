package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type DirectMessageReactionId string

type DirectMessageReaction struct {
	Id       DirectMessageReactionId
	UserIds  []user_entity.UserId
	Reaction string
}

func (id DirectMessageReactionId) String() string {
	return string(id)
}
