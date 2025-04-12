package outbound

import (
	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundUserStatusUpdated struct {
	messages.DefaultMessage
	UserId user_entity.UserId `json:"userId"`
	Status entity.Status      `json:"status"`
}

func (m OutboundUserStatusUpdated) GetActionName() messages.Action {
	return messages.OutboundUserStatusUpdatedAction
}

func (m OutboundUserStatusUpdated) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
