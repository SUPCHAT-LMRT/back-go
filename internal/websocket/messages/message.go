package messages

import (
	"github.com/google/uuid"
	"time"
)

type Action string

// Outbound actions are actions that are sent from the server to the client.
const (
	OutboundSendChannelMessageAction     Action = "send-channel-message"
	OutboundRoomJoinedAction             Action = "room-joined"
	OutboundChannelCreatedAction         Action = "channel-created"
	OutboundChannelMessageReactionCreate Action = "channel-message-reaction-create"
)

// Inbound actions are actions that are sent from the client to the server.
const (
	InboundSendChannelMessageAction     Action = "send-channel-message"
	InboundJoinDirectRoomAction         Action = "join-direct-room"
	InboundJoinGroupRoomAction          Action = "join-group-room"
	InboundJoinChannelRoomAction        Action = "join-channel-room"
	InboundLeaveRoomAction              Action = "leave-room"
	InboundUnselectWorkspaceAction      Action = "unselect-workspace"
	InboundSelectWorkspaceAction        Action = "select-workspace"
	InboundChannelMessageReactionCreate Action = "channel-message-reaction-create"
)

type Message interface {
	GetActionName() Action
	SetId(string)
	SetCreatedAt(time.Time)
	Encode() ([]byte, error)
	mustExtendDefaultMessage()
}

type DefaultMessage struct {
	Id        string    `json:"id"`
	Action    Action    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *DefaultMessage) SetId(id string) {
	m.Id = id
}

func (m *DefaultMessage) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *DefaultMessage) mustExtendDefaultMessage() {}

func NewDefaultMessage(action Action) DefaultMessage {
	return DefaultMessage{Id: uuid.NewString(), Action: action, CreatedAt: time.Now()}
}
