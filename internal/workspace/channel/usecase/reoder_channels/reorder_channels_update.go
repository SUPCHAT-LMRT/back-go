package reoder_channels

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type ReorderChannelsObserverDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type ReorderChannelsObserver struct {
	deps ReorderChannelsObserverDeps
}

func NewUserStatusUpdateObserver(deps ReorderChannelsObserverDeps) ReorderIndexChannelsObserver {
	return &ReorderChannelsObserver{deps: deps}
}

func (o ReorderChannelsObserver) NotifyChannelReordered(channels []ChannelReorderMessage, workspaceId entity.WorkspaceId) {

	var eventChannelReorders []event.ChannelReorderMessage
	for _, channel := range channels {
		eventChannelReorders = append(eventChannelReorders, event.ChannelReorderMessage{
			ChannelId: channel.ChannelId,

			NewOrder: channel.NewOrder,
		})
	}

	o.deps.EventBus.Publish(&event.ChannelsReorderedEvent{
		ChannelReorders: eventChannelReorders,
		WorkspaceId:     workspaceId,
	})
}
