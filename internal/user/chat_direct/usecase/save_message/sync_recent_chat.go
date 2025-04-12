package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
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

func NewSyncRecentChatObserver(deps SyncRecentChatObserverDeps) MessageSavedObserver {
	return &SyncRecentChatObserver{deps: deps}
}

func (o SyncRecentChatObserver) NotifyMessageSaved(msg *entity.ChatDirect) {
	// Check if the message is the first message between the two users
	isFirst, err := o.deps.IsFirstMessageUseCase.Execute(context.Background(), msg.User1Id, msg.User2Id)
	if err != nil {
		return
	}

	if !isFirst {
		return
	}

	// Publish an event after saving the message
	o.deps.EventBus.Publish(&event.DirectChatMessageSavedEvent{
		Message: msg,
	})
}
