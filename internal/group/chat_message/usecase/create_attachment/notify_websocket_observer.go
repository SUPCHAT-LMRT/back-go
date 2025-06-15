package create_attachment

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/logger"
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

func (o NotifyWebsocketObserver) NotifyAttachmentCreated(message *entity.GroupChatMessage) {
	o.deps.EventBus.Publish(&event.GroupAttachmentSentEvent{
		GroupChatMessage: message,
	})
}
