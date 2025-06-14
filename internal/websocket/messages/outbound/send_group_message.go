package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"time"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundSendGroupMessageSender struct {
	UserId    user_entity.UserId `json:"userId"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

type OutboundSendGroupMessage struct {
	messages.DefaultMessage
	MessageId string                          `json:"messageId"`
	Content   string                          `json:"content"`
	GroupId   string                          `json:"groupId"`
	Sender    *OutboundSendGroupMessageSender `json:"sender"`
	CreatedAt time.Time                       `json:"createdAt"`
}

func (m *OutboundSendGroupMessage) GetActionName() messages.Action {
	return messages.OutboundSendGroupMessageAction
}

func (m *OutboundSendGroupMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

type OutboundGroupMessageReactionMember struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type OutboundGroupMessageReactionAdded struct {
	messages.DefaultMessage
	GroupId   group_entity.GroupId               `json:"groupId"`
	MessageId string                             `json:"messageId"`
	Reaction  string                             `json:"reaction"`
	Member    OutboundGroupMessageReactionMember `json:"member"`
}

func (m *OutboundGroupMessageReactionAdded) GetActionName() messages.Action {
	return messages.OutboundGroupMessageReactionAddedAction
}

func (m *OutboundGroupMessageReactionAdded) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

type OutboundGroupMessageReactionRemoved struct {
	messages.DefaultMessage
	MessageId string                             `json:"messageId"`
	Reaction  string                             `json:"reaction"`
	GroupId   group_entity.GroupId               `json:"groupId"`
	Member    OutboundGroupMessageReactionMember `json:"member"`
}

func (m *OutboundGroupMessageReactionRemoved) GetActionName() messages.Action {
	return messages.OutboundGroupMessageReactionRemovedAction
}

func (m *OutboundGroupMessageReactionRemoved) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

type OutboundGroupRoomJoined struct {
	messages.DefaultMessage
	RoomId string `json:"roomId"`
}

func (m *OutboundGroupRoomJoined) GetActionName() messages.Action {
	return messages.OutboundGroupRoomJoinedAction
}

func (m *OutboundGroupRoomJoined) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
