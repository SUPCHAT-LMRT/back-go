package outbound

import (
	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundDirectRoomJoined struct {
	messages.DefaultMessage
	RoomId      string             `json:"roomId"`
	OtherUserId user_entity.UserId `json:"otherUserId"`
}

func (m *OutboundDirectRoomJoined) GetActionName() messages.Action {
	return messages.OutboundDirectRoomJoinedAction
}

func (m *OutboundDirectRoomJoined) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
