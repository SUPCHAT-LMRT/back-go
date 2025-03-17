package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundChannelRoomJoined struct {
	messages.DefaultMessage
	RoomId string `json:"roomId"`
}

func (m OutboundChannelRoomJoined) GetActionName() messages.Action {
	return messages.OutboundChannelRoomJoinedAction
}

func (m OutboundChannelRoomJoined) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
