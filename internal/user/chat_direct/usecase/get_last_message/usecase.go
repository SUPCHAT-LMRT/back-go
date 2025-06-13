package get_last_message

import (
	"context"
	"time"

	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type GetLastDirectChatMessageUseCaseDeps struct {
	uberdig.In
	ChatDirectRepository repository.ChatDirectRepository
}

type GetLastDirectChatMessageUseCase struct {
	deps GetLastDirectChatMessageUseCaseDeps
}

func NewGetLastDirectChatMessageUseCase(
	deps GetLastDirectChatMessageUseCaseDeps,
) *GetLastDirectChatMessageUseCase {
	return &GetLastDirectChatMessageUseCase{deps: deps}
}

func (uc *GetLastDirectChatMessageUseCase) Execute(
	ctx context.Context,
	user1Id, user2Id user_entity.UserId,
) (*LastMessageResponse, error) {
	lastMessage, err := uc.deps.ChatDirectRepository.GetLastMessage(ctx, user1Id, user2Id)
	if err != nil {
		return nil, err
	}

	if lastMessage == nil {
		return nil, nil // No messages found
	}

	return &LastMessageResponse{
		ChatId:    lastMessage.Id,
		Content:   lastMessage.Content,
		CreatedAt: lastMessage.CreatedAt,
		AuthorId:  lastMessage.SenderId,
	}, nil
}

type LastMessageResponse struct {
	ChatId    chat_direct_entity.ChatDirectId
	Content   string
	CreatedAt time.Time
	AuthorId  user_entity.UserId
}
