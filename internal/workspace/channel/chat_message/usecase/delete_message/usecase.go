package delete_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	uberdig "go.uber.org/dig"
)

type DeleteChannelChatMessageUseCaseDeps struct {
	uberdig.In
	ChannelMessageRepository repository2.ChannelMessageRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type DeleteChannelChatMessageUseCase struct {
	deps DeleteChannelChatMessageUseCaseDeps
}

func NewDeleteChannelChatMessageUseCase(
	deps DeleteChannelChatMessageUseCaseDeps,
) *DeleteChannelChatMessageUseCase {
	return &DeleteChannelChatMessageUseCase{deps: deps}
}

func (u *DeleteChannelChatMessageUseCase) Execute(
	ctx context.Context,
	msgId channel_message_entity.ChannelMessageId,
) error {
	if err := u.deps.ChannelMessageRepository.DeleteMessage(ctx, msgId); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.DeleteMessage(ctx, msgId.String())
	if err != nil {
		return err
	}

	return nil
}
