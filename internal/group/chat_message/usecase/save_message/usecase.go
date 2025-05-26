package save_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
)

type SaveGroupChatMessageUseCase struct {
	repository repository.GroupChatMessageRepository
}

func NewSaveGroupChatMessageUseCase(
	repository repository.GroupChatMessageRepository,
) *SaveGroupChatMessageUseCase {
	return &SaveGroupChatMessageUseCase{repository: repository}
}

func (u SaveGroupChatMessageUseCase) Execute(
	ctx context.Context,
	message *entity.GroupChatMessage,
) error {
	return u.repository.Create(ctx, message)
}
