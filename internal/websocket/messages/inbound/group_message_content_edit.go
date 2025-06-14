package inbound

import (
	"github.com/goccy/go-json"
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundGroupMessageContentEdit struct {
	messages.DefaultMessage
	GroupId    group_entity.GroupId                 `json:"groupId"`
	MessageId  group_chat_entity.GroupChatMessageId `json:"messageId"`
	NewContent string                               `json:"newContent"`
}

func (i *InboundGroupMessageContentEdit) GetActionName() messages.Action {
	return messages.InboundGroupMessageContentEdit
}

func (i *InboundGroupMessageContentEdit) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
