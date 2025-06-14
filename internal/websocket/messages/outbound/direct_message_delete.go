package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundDirectMessageDeleted struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId  `json:"otherUserId"`
	MessageId   entity.ChatDirectId `json:"messageId"`
}

func (o *OutboundDirectMessageDeleted) GetActionName() messages.Action {
	return messages.OutboundDirectMessageDeletedAction
}

func (o *OutboundDirectMessageDeleted) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
