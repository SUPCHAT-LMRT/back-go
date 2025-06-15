package create_attachment

import (
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type CreateAttachmentObserver interface {
	NotifyAttachmentCreated(workspaceId workspace_entity.WorkspaceId, workspaceMemberId workspace_member_entity.WorkspaceMemberId, message *channel_message_entity.ChannelMessage)
}
