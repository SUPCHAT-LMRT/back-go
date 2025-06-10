package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundChannelMessageReactionToggle struct {
	messages.DefaultMessage
	RoomId    string `json:"roomId"`
	MessageId string `json:"messageId"`
	Reaction  string `json:"reaction"`
}

func (i *InboundChannelMessageReactionToggle) GetActionName() messages.Action {
	return messages.InboundChannelMessageReactionToggle
}

func (i *InboundChannelMessageReactionToggle) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
