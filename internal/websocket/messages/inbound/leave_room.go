package inbound

import (
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundLeaveRoom struct {
	messages.DefaultMessage
	RoomId string `json:"roomId"`
}

func (m InboundLeaveRoom) GetActionName() messages.Action {
	return messages.InboundJoinChannelRoomAction
}
