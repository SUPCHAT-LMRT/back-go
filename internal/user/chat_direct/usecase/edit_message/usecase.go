package edit_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	uberdig "go.uber.org/dig"
	"time"
)

type EditDirectChatMessageUseCaseDeps struct {
	uberdig.In
	ChatDirectRepository     repository.ChatDirectRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type EditDirectChatMessageUseCase struct {
	deps EditDirectChatMessageUseCaseDeps
}

func NewEditDirectChatMessageUseCase(
	deps EditDirectChatMessageUseCaseDeps,
) *EditDirectChatMessageUseCase {
	return &EditDirectChatMessageUseCase{deps: deps}
}

func (u *EditDirectChatMessageUseCase) Execute(
	ctx context.Context,
	msg *chat_direct_entity.ChatDirect,
) error {
	if msg.User2Id.IsAfter(msg.User1Id) {
		msg.User1Id, msg.User2Id = msg.User2Id, msg.User1Id
	}

	msg.UpdatedAt = time.Now()
	if err := u.deps.ChatDirectRepository.UpdateMessage(ctx, msg); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageGroupMessage,
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

	return nil
}
