package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundDirectMessageDelete struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId  `json:"otherUserId"`
	MessageId   entity.ChatDirectId `json:"messageId"`
}

func (i *InboundDirectMessageDelete) GetActionName() messages.Action {
	return messages.InboundDirectMessageDelete
}

func (i *InboundDirectMessageDelete) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
