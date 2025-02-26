package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundMessageReactionCreate struct {
	messages.DefaultMessage
	RoomId    string `json:"roomId"`
	MessageId string `json:"messageId"`
	Reaction  string `json:"reaction"`
}

func (i InboundMessageReactionCreate) GetActionName() messages.Action {
	return messages.InboundChannelMessageReactionCreate
}

func (i InboundMessageReactionCreate) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
