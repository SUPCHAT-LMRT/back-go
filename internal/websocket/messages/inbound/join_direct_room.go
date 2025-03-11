package inbound

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundJoinDirectRoom struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId `json:"otherUserId"`
}

func (m InboundJoinDirectRoom) GetActionName() messages.Action {
	return messages.InboundJoinDirectRoomAction
}
