package outbound

import (
	"github.com/goccy/go-json"
	"time"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundSendGroupMessageSender struct {
	UserId    user_entity.UserId `json:"user_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
}

type OutboundSendGroupMessage struct {
	messages.DefaultMessage
	MessageId string                          `json:"message_id"`
	Content   string                          `json:"content"`
	GroupId   string                          `json:"group_id"`
	Sender    *OutboundSendGroupMessageSender `json:"sender"`
	CreatedAt time.Time                       `json:"created_at"`
}

func (m *OutboundSendGroupMessage) GetActionName() messages.Action {
	return messages.OutboundSendGroupMessageAction
}

func (m *OutboundSendGroupMessage) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

type OutboundGroupMessageReactionMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OutboundGroupMessageReactionAdded struct {
	messages.DefaultMessage
	MessageId string                             `json:"message_id"`
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
	MessageId string                             `json:"message_id"`
	Reaction  string                             `json:"reaction"`
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
	RoomId  string `json:"room_id"`
	GroupId string `json:"group_id"`
}

func (m *OutboundGroupRoomJoined) GetActionName() messages.Action {
	return messages.OutboundGroupRoomJoinedAction
}

func (m *OutboundGroupRoomJoined) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
