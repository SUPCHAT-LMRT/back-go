package inbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundGroupMessageReactionToggle struct {
	messages.DefaultMessage
	GroupId   group_entity.GroupId `json:"groupId"`
	MessageId string               `json:"messageId"`
	Reaction  string               `json:"reaction"`
}

func (i *InboundGroupMessageReactionToggle) GetActionName() messages.Action {
	return messages.InboundDirectMessageReactionToggle
}

func (i *InboundGroupMessageReactionToggle) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
