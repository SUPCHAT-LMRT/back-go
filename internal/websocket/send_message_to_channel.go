package websocket

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SendMessageToChannel struct {
	DefaultMessage
	Sender    *SendMessageToChannelSender `json:"sender"`
	Content   string                      `json:"content"`
	ChannelId channel_entity.ChannelId    `json:"channelId"`
}

type SendMessageToChannelSender struct {
	UserId            user_entity.UserId       `json:"userId"`
	Pseudo            string                   `json:"pseudo"`
	WorkspaceMemberId entity.WorkspaceMemberId `json:"workspaceMemberId"`
	WorkspacePseudo   string                   `json:"workspacePseudo"`
}

func (m SendMessageToChannel) GetActionName() string {
	return OutboundSendMessageAction
}
