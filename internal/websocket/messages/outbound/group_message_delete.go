package outbound

import (
	"github.com/goccy/go-json"
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundGroupMessageDeleted struct {
	messages.DefaultMessage
	GroupId   group_entity.GroupId                 `json:"groupId"`
	MessageId group_chat_entity.GroupChatMessageId `json:"messageId"`
}

func (o *OutboundGroupMessageDeleted) GetActionName() messages.Action {
	return messages.OutboundGroupMessageDeletedAction
}

func (o *OutboundGroupMessageDeleted) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
