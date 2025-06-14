package edit_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	uberdig "go.uber.org/dig"
	"time"
)

type EditGroupChatMessageUseCaseDeps struct {
	uberdig.In
	ChatMessageRepository    repository.ChatMessageRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type EditGroupChatMessageUseCase struct {
	deps EditGroupChatMessageUseCaseDeps
}

func NewEditGroupChatMessageUseCase(
	deps EditGroupChatMessageUseCaseDeps,
) *EditGroupChatMessageUseCase {
	return &EditGroupChatMessageUseCase{deps: deps}
}

func (u *EditGroupChatMessageUseCase) Execute(
	ctx context.Context,
	msg *entity.GroupChatMessage,
) error {
	msg.UpdatedAt = time.Now()
	if err := u.deps.ChatMessageRepository.UpdateMessage(ctx, msg); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageGroupMessage,
		Data: message.SearchMessageGroupData{
			GroupId: msg.GroupId,
		},
		AuthorId:  msg.AuthorId,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
