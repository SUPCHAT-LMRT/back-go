package websocket

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type SendMessageObserver interface {
	OnSendMessage(message messages.Message, userId user_entity.UserId)
}
