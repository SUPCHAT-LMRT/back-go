package websocket

import (
	"github.com/goccy/go-json"
	"log"
	"time"
)

// Outbound actions are actions that are sent from the server to the client.
const OutboundSendMessageAction = "send-message"
const OutboundRoomJoinedAction = "room-joined"
const OutboundChannelCreatedAction = "channel-created"

// Inbound actions are actions that are sent from the client to the server.
const InboundJoinDirectRoomAction = "join-direct-room"
const InboundJoinGroupRoomAction = "join-group-room"
const InboundJoinChannelRoomAction = "join-channel-room"
const InboundLeaveRoomAction = "leave-room"
const InboundUnselectWorkspaceAction = "unselect-workspace"
const InboundSelectWorkspaceAction = "select-workspace"

type Message interface {
	GetActionName() string
	SetId(string)
	SetCreatedAt(time.Time)
	encode() []byte
	SetEmittedBy(*Client)
}

type DefaultMessage struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	EmittedBy *Client   `json:"emittedBy"`
}

func (m *DefaultMessage) SetId(id string) {
	m.Id = id
}

func (m *DefaultMessage) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *DefaultMessage) SetEmittedBy(client *Client) {
	m.EmittedBy = client
}

func (m *DefaultMessage) encode() []byte {
	result, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}

	return result
}
