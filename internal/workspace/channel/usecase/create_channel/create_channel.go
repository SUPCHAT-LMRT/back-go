package create_channel

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	uberdig "go.uber.org/dig"
)

type CreateChannelObserverDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type CreateChannelObserver struct {
	deps CreateChannelObserverDeps
}

func NewCreateChannelObserver(deps CreateChannelObserverDeps) CreateSpecifyChannelObserver {
	return &CreateChannelObserver{deps: deps}
}

func (o CreateChannelObserver) NotifyChannelCreated(channel *entity.Channel) {
	o.deps.EventBus.Publish(&event.ChannelCreatedEvent{
		Channel: channel,
	})
}
