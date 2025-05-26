package list_recent_direct_chats

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type ListRecentChatDirectUseCaseDeps struct {
	uberdig.In
	Repository repository.ChatDirectRepository
}

type ListRecentChatDirectUseCase struct {
	deps ListRecentChatDirectUseCaseDeps
}

func NewListRecentChatDirectUseCase(
	deps ListRecentChatDirectUseCaseDeps,
) *ListRecentChatDirectUseCase {
	return &ListRecentChatDirectUseCase{deps: deps}
}

func (u *ListRecentChatDirectUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
) ([]*entity.ChatDirect, error) {
	return u.deps.Repository.ListRecentChats(ctx, userId)
}
