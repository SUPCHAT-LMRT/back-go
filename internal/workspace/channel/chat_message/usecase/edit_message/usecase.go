package edit_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	uberdig "go.uber.org/dig"
	"time"
)

type EditChannelChatMessageUseCaseDeps struct {
	uberdig.In
	ChannelMessageRepository repository2.ChannelMessageRepository
	GetChannelUseCase        *get_channel.GetChannelUseCase
	SearchMessageSyncManager message.SearchMessageSyncManager
}

type EditChannelChatMessageUseCase struct {
	deps EditChannelChatMessageUseCaseDeps
}

func NewEditChannelChatMessageUseCase(
	deps EditChannelChatMessageUseCaseDeps,
) *EditChannelChatMessageUseCase {
	return &EditChannelChatMessageUseCase{deps: deps}
}

func (u *EditChannelChatMessageUseCase) Execute(
	ctx context.Context,
	msg *entity2.ChannelMessage,
) error {
	msg.UpdatedAt = time.Now()
	if err := u.deps.ChannelMessageRepository.UpdateMessage(ctx, msg); err != nil {
		return err
	}

	channel, err := u.deps.GetChannelUseCase.Execute(ctx, msg.ChannelId)
	if err != nil {
		return err
	}

	err = u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageGroupMessage,
		Data: message.SearchMessageChannelData{
			ChannelId:   msg.ChannelId,
			WorkspaceId: channel.WorkspaceId,
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
