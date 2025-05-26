package delete_channels

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type DeleteChannelsObserverDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type DeleteChannelsObserver struct {
	deps DeleteChannelsObserverDeps
}

func NewDeleteChannelsObserver(deps DeleteChannelsObserverDeps) DeleteSpecifyChannelsObserver {
	return &DeleteChannelsObserver{deps: deps}
}

func (o DeleteChannelsObserver) NotifyChannelsDeleted(
	channelId entity.ChannelId,
	workspaceId workspace_entity.WorkspaceId,
) {
	o.deps.EventBus.Publish(&event.ChannelsDeletedEvent{
		ChannelId:   channelId,
		WorkspaceId: workspaceId,
	})
}
