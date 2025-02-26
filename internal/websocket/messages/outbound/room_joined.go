package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/websocket/room"
)

type OutboundRoomJoined struct {
	messages.DefaultMessage
	Room OutboundRoomJoinedRoom `json:"room"`
}

type OutboundRoomJoinedRoom struct {
	Id   string        `json:"id"`
	Kind room.RoomKind `json:"kind"`
}

func (m OutboundRoomJoined) GetActionName() messages.Action {
	return messages.OutboundRoomJoinedAction
}

func (m OutboundRoomJoined) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
