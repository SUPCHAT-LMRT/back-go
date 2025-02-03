package entity

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type ChatDirectId string

type ChatDirect struct {
	Id        ChatDirectId
	User1Id   user_entity.UserId
	User2Id   user_entity.UserId
	CreatedAt time.Time
	UpdatedAt time.Time
}
