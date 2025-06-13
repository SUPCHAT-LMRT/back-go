package inbound

import (
	"github.com/goccy/go-json"
	"time"

	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundSendGroupMessage struct {
	messages.DefaultMessage
	GroupId                   string    `json:"group_id"`
	Content                   string    `json:"content"`
	TransportMessageCreatedAt time.Time `json:"created_at"`
}

func (m *InboundSendGroupMessage) GetActionName() messages.Action {
	return messages.InboundSendGroupMessageAction
}

func (m *InboundSendGroupMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
