package save_message

import (
	"context"
	chat_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
)

type SaveChannelMessageUseCase struct {
	repository chat_message_repository.ChannelMessageRepository
}

func NewSaveChannelMessageUseCase(repository chat_message_repository.ChannelMessageRepository) *SaveChannelMessageUseCase {
	return &SaveChannelMessageUseCase{repository: repository}
}

func (u SaveChannelMessageUseCase) Execute(ctx context.Context, message *chat_message_entity.ChannelMessage) error {
	return u.repository.Create(ctx, message)
}
