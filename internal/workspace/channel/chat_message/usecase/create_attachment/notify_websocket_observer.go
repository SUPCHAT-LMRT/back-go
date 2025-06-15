package create_attachment

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	uberdig "go.uber.org/dig"
)

type NotifyWebsocketObserverDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type NotifyWebsocketObserver struct {
	deps NotifyWebsocketObserverDeps
}

func NewNotifyWebsocketObserver(deps NotifyWebsocketObserverDeps) CreateAttachmentObserver {
	return &NotifyWebsocketObserver{deps: deps}
}

func (o NotifyWebsocketObserver) NotifyAttachmentCreated(workspaceId workspace_entity.WorkspaceId, workspaceMemberId workspace_member_entity.WorkspaceMemberId, message *channel_message_entity.ChannelMessage) {
	o.deps.EventBus.Publish(&event.ChannelAttachmentSentEvent{
		ChannelMessage:    message,
		WorkspaceId:       workspaceId,
		WorkspaceMemberId: workspaceMemberId,
	})
}
