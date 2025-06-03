package outbound

import (
	"time"

	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundSendDirectMessage struct {
	messages.DefaultMessage
	MessageId   string                           `json:"messageId"`
	Sender      *OutboundSendDirectMessageSender `json:"sender"`
	Content     string                           `json:"content"`
	OtherUserId user_entity.UserId               `json:"otherUserId"`
	CreatedAt   time.Time                        `json:"createdAt"`
}

type OutboundSendDirectMessageSender struct {
	UserId    user_entity.UserId `json:"userId"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

func (m OutboundSendDirectMessage) GetActionName() messages.Action {
	return messages.OutboundSendDirectMessageAction
}

func (m OutboundSendDirectMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
