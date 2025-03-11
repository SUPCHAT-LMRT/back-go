package save_message

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	uberdig "go.uber.org/dig"
)

type SaveDirectMessageUseCaseDeps struct {
	uberdig.In
	Repository repository.ChatDirectRepository
}

type SaveDirectMessageUseCase struct {
	deps SaveDirectMessageUseCaseDeps
}

func NewSaveDirectMessageUseCase(deps SaveDirectMessageUseCaseDeps) *SaveDirectMessageUseCase {
	return &SaveDirectMessageUseCase{deps: deps}
}

func (u SaveDirectMessageUseCase) Execute(ctx context.Context, message *chat_direct_entity.ChatDirect) error {
	// Ensure that the two users are correctly ordered (the one with the smallest ID is the first)
	if message.User2Id.IsAfter(message.User1Id) {
		message.User1Id, message.User2Id = message.User2Id, message.User1Id
	}
	return u.deps.Repository.Create(ctx, message)
}
