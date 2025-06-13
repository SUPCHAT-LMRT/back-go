package add_member

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

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

func NewSyncRecentChatObserver(deps SyncRecentChatObserverDeps) AddGroupMemberObserver {
	return &SyncRecentChatObserver{deps: deps}
}

func (o SyncRecentChatObserver) NotifyGroupMemberKicked(msg *group_entity.Group, inviterUserId user_entity.UserId) {
	// Publish an event after adding a member to the group
	o.deps.EventBus.Publish(&event.GroupMemberAddedEvent{
		Message:       msg,
		InvitedUserId: inviterUserId,
	})
}
