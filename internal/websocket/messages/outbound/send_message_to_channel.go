package outbound

import (
	"time"

	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type OutboundSendMessageToChannel struct {
	messages.DefaultMessage
	MessageId string                              `json:"messageId"`
	Sender    *OutboundSendMessageToChannelSender `json:"sender"`
	Content   string                              `json:"content"`
	ChannelId channel_entity.ChannelId            `json:"channelId"`
	CreatedAt time.Time                           `json:"createdAt"`
}

type OutboundSendMessageToChannelSender struct {
	UserId            user_entity.UserId                 `json:"userId"`
	Pseudo            string                             `json:"pseudo"`
	WorkspaceMemberId workspace_entity.WorkspaceMemberId `json:"workspaceMemberId"`
	WorkspacePseudo   string                             `json:"workspacePseudo"`
}

func (m *OutboundSendMessageToChannel) GetActionName() messages.Action {
	return messages.OutboundSendChannelMessageAction
}

func (m *OutboundSendMessageToChannel) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
