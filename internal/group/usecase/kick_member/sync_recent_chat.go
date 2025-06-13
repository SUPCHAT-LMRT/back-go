package kick_member

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/is_first_message"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
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

func NewSyncRecentChatObserver(deps SyncRecentChatObserverDeps) KickGroupMemberObserver {
	return &SyncRecentChatObserver{deps: deps}
}

func (o SyncRecentChatObserver) NotifyGroupMemberKicked(msg *group_entity.Group, kickedMemberId group_entity.GroupMemberId, kickedUserId user_entity.UserId) {
	// Publish an event after kicking a member from the group
	o.deps.EventBus.Publish(&event.GroupMemberRemovedEvent{
		Group:           msg,
		RemovedMemberId: kickedMemberId,
		RemovedUserId:   kickedUserId,
	})
}
