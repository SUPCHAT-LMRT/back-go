package create_attachment

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
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

func (o NotifyWebsocketObserver) NotifyAttachmentCreated(message *chat_direct_entity.ChatDirect) {
	o.deps.EventBus.Publish(&event.ChatDirectAttachmentSentEvent{
		ChatDirect: message,
	})
}
