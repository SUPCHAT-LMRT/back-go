package inbound

import (
	"github.com/goccy/go-json"
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

func (m InboundJoinDirectRoom) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
