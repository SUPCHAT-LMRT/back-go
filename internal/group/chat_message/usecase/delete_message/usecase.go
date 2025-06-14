package delete_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	uberdig "go.uber.org/dig"
)

type DeleteGroupChatMessageUseCaseDeps struct {
	uberdig.In
	ChatMessageRepository    repository.ChatMessageRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type DeleteGroupChatMessageUseCase struct {
	deps DeleteGroupChatMessageUseCaseDeps
}

func NewDeleteGroupChatMessageUseCase(
	deps DeleteGroupChatMessageUseCaseDeps,
) *DeleteGroupChatMessageUseCase {
	return &DeleteGroupChatMessageUseCase{deps: deps}
}

func (u *DeleteGroupChatMessageUseCase) Execute(
	ctx context.Context,
	msgId entity.GroupChatMessageId,
) error {
	if err := u.deps.ChatMessageRepository.DeleteMessage(ctx, msgId); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.DeleteMessage(ctx, msgId.String())
	if err != nil {
		return err
	}

	return nil
}
