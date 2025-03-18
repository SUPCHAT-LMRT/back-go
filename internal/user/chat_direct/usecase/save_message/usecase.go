package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/is_first_message"
	uberdig "go.uber.org/dig"
)

type SaveDirectMessageUseCaseDeps struct {
	uberdig.In
	Repository               repository.ChatDirectRepository
	IsFirstMessageUseCase    *is_first_message.IsFirstMessageUseCase
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type SaveDirectMessageUseCase struct {
	deps SaveDirectMessageUseCaseDeps
}

func NewSaveDirectMessageUseCase(deps SaveDirectMessageUseCaseDeps) *SaveDirectMessageUseCase {
	return &SaveDirectMessageUseCase{deps: deps}
}

func (u SaveDirectMessageUseCase) Execute(ctx context.Context, msg *chat_direct_entity.ChatDirect) error {
	// Ensure that the two users are correctly ordered (the one with the smallest ID is the first)
	if msg.User2Id.IsAfter(msg.User1Id) {
		msg.User1Id, msg.User2Id = msg.User2Id, msg.User1Id
	}
	err := u.deps.Repository.Create(ctx, msg)
	if err != nil {
		return err
	}

	err = u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageKindDirectMessage,
		Data: message.SearchMessageDirectData{
			User1: msg.User1Id,
			User2: msg.User2Id,
		},
		AuthorId:  msg.SenderId,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return err
}
