package delete_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	uberdig "go.uber.org/dig"
)

type DeleteDirectChatMessageUseCaseDeps struct {
	uberdig.In
	ChatDirectRepository     repository.ChatDirectRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type DeleteDirectChatMessageUseCase struct {
	deps DeleteDirectChatMessageUseCaseDeps
}

func NewDeleteDirectChatMessageUseCase(
	deps DeleteDirectChatMessageUseCaseDeps,
) *DeleteDirectChatMessageUseCase {
	return &DeleteDirectChatMessageUseCase{deps: deps}
}

func (u *DeleteDirectChatMessageUseCase) Execute(
	ctx context.Context,
	msgId chat_direct_entity.ChatDirectId,
) error {
	if err := u.deps.ChatDirectRepository.DeleteMessage(ctx, msgId); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.DeleteMessage(ctx, msgId.String())
	if err != nil {
		return err
	}

	return nil
}
