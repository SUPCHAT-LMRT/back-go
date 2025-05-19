package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	uberdig "go.uber.org/dig"
)

type SaveDirectMessageUseCaseDeps struct {
	uberdig.In
	Repository               repository.ChatDirectRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
	Observers                []MessageSavedObserver `group:"save_direct_chat_message_observers"`
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

	// donc en gros pour faire gros, j'appelle mon usecase SendNotificationUseCase ici
	for _, observer := range u.deps.Observers {
		observer.NotifyMessageSaved(msg)
	}

	return err
}
