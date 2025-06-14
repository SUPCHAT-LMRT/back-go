package create_group

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"

	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/is_first_message"
	uberdig "go.uber.org/dig"
)

type SyncRecentChatObserverDeps struct {
	uberdig.In
	IsFirstMessageUseCase *is_first_message.IsFirstMessageUseCase
	EventBus              *event.EventBus
	Logger                logger.Logger
}

type SyncRecentChatObserver struct {
	deps SyncRecentChatObserverDeps
}

func NewSyncRecentChatObserver(deps SyncRecentChatObserverDeps) GroupCreatedObserver {
	return &SyncRecentChatObserver{deps: deps}
}

func (o SyncRecentChatObserver) NotifyGroupMemberAdded(msg *group_entity.Group) {
	// Publish an event after creating the group
	o.deps.EventBus.Publish(&event.GroupCreatedEvent{
		Message: msg,
	})
}
