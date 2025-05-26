package inbound

import (
	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundSendDirectMessage struct {
	messages.DefaultMessage
	Content     string             `json:"content"`
	OtherUserId user_entity.UserId `json:"otherUserId"`
}

func (m *InboundSendDirectMessage) GetActionName() messages.Action {
	return messages.InboundSendDirectMessageAction
}

func (m *InboundSendDirectMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
