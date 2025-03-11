package list_messages

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ListDirectMessagesUseCase struct {
	repository repository.ChatDirectRepository
}

func NewListDirectMessagesUseCase(repository repository.ChatDirectRepository) *ListDirectMessagesUseCase {
	return &ListDirectMessagesUseCase{repository: repository}
}

func (u ListDirectMessagesUseCase) Execute(ctx context.Context, user1Id, user2Id user_entity.UserId) ([]*chat_direct_entity.ChatDirect, error) {
	return u.repository.ListByUser(ctx, user1Id, user2Id)
}
