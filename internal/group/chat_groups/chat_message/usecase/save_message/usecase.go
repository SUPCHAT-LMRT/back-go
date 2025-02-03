package save_message

import (
	"context"
	chat_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
)

type SaveGroupChatMessageUseCase struct {
	repository chat_message_repository.ChannelMessageRepository
}

func NewSaveGroupChatMessageUseCase(repository chat_message_repository.ChannelMessageRepository) *SaveGroupChatMessageUseCase {
	return &SaveGroupChatMessageUseCase{repository: repository}
}

func (u SaveGroupChatMessageUseCase) Execute(ctx context.Context, message *chat_message_entity.ChannelMessage) error {
	return u.repository.Create(ctx, message)
}
