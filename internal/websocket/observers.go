package websocket

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
)

type SendMessageObserver interface {
	OnSendMessage(message messages.Message, messageId entity.ChannelMessageId, userId user_entity.UserId)
}
