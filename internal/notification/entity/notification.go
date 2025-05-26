package entity

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type NotificationId string

type Notification struct {
	Id        NotificationId
	UserId    entity.UserId
	Content   string
	IsRead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (id NotificationId) String() string {
	return string(id)
}
