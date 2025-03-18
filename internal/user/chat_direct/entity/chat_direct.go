package entity

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type ChatDirectId string

type ChatDirect struct {
	Id        ChatDirectId
	SenderId  user_entity.UserId
	User1Id   user_entity.UserId
	User2Id   user_entity.UserId
	Content   string
	Reactions []*DirectMessageReaction
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (id ChatDirectId) String() string {
	return string(id)
}
