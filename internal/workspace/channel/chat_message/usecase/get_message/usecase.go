package get_message

import (
	"context"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
)

type GetMessageUseCase struct {
	repository repository.ChannelMessageRepository
}

func NewGetMessageUseCase(repository repository.ChannelMessageRepository) *GetMessageUseCase {
	return &GetMessageUseCase{
		repository: repository,
	}
}

func (uc *GetMessageUseCase) Execute(ctx context.Context, messageId channel_message_entity.ChannelMessageId) (*channel_message_entity.ChannelMessage, error) {
	return uc.repository.Get(ctx, messageId)
}
