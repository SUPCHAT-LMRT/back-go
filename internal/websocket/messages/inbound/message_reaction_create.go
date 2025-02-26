package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundMessageReactionToggle struct {
	messages.DefaultMessage
	RoomId    string `json:"roomId"`
	MessageId string `json:"messageId"`
	Reaction  string `json:"reaction"`
}

func (i InboundMessageReactionToggle) GetActionName() messages.Action {
	return messages.InboundChannelMessageReactionToggle
}

func (i InboundMessageReactionToggle) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
