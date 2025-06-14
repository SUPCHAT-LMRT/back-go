package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundGroupOwnershipTransferer struct {
	messages.DefaultMessage
	GroupId    group_entity.GroupId       `json:"groupId"`
	NewOwnerId group_entity.GroupMemberId `json:"newOwnerId"`
}

func (o *OutboundGroupOwnershipTransferer) GetActionName() messages.Action {
	return messages.OutboundGroupOwnershipTransferredAction
}

func (o *OutboundGroupOwnershipTransferer) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
