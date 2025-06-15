package event

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

const (
	ChannelAttachmentSentEventType EventType = "channel_attachment_sent"
)

type ChannelAttachmentSentEvent struct {
	ChannelMessage    *entity.ChannelMessage
	WorkspaceId       workspace_entity.WorkspaceId
	WorkspaceMemberId workspace_member_entity.WorkspaceMemberId
}

func (e ChannelAttachmentSentEvent) Type() EventType {
	return ChannelAttachmentSentEventType
}
