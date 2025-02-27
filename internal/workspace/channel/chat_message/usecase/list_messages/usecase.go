package list_messages

import (
	"context"
	chat_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ListChannelMessagesUseCase struct {
	repository chat_message_repository.ChannelMessageRepository
}

func NewListMessageUseCase(repository chat_message_repository.ChannelMessageRepository) *ListChannelMessagesUseCase {
	return &ListChannelMessagesUseCase{repository: repository}
}

func (u ListChannelMessagesUseCase) Execute(ctx context.Context, channelId entity.ChannelId) ([]*chat_message_entity.ChannelMessage, error) {
	return u.repository.ListByChannelId(ctx, channelId)
}
