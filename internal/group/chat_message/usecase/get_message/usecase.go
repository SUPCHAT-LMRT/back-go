package get_message

import (
	"context"
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
)

type GetMessageUseCase struct {
	repository repository.ChatMessageRepository
}

func NewGetMessageUseCase(repository repository.ChatMessageRepository) *GetMessageUseCase {
	return &GetMessageUseCase{
		repository: repository,
	}
}

func (uc *GetMessageUseCase) Execute(ctx context.Context, messageId group_chat_entity.GroupChatMessageId) (*group_chat_entity.GroupChatMessage, error) {
	return uc.repository.Get(ctx, messageId)
}
