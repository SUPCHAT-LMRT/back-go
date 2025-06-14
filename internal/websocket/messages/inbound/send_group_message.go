package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundSendGroupMessage struct {
	messages.DefaultMessage
	GroupId string `json:"groupId"`
	Content string `json:"content"`
}

func (m *InboundSendGroupMessage) GetActionName() messages.Action {
	return messages.InboundSendGroupMessageAction
}

func (m *InboundSendGroupMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
