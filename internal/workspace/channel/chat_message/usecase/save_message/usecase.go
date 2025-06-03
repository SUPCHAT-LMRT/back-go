package save_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/search/message"
	chat_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/entity"
	time_series_message_sent_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type SaveChannelMessageUseCaseDeps struct {
	uberdig.In
	Repository                chat_message_repository.ChannelMessageRepository
	TimeSeriesRepository      time_series_message_sent_repository.MessageSentTimeSeriesWorkspaceRepository
	GetChannelUseCase         *get_channel.GetChannelUseCase
	GetWorkspaceMemberUseCase *get_workpace_member.GetWorkspaceMemberUseCase
	SearchMessageSyncManager  message.SearchMessageSyncManager
}

type SaveChannelMessageUseCase struct {
	deps SaveChannelMessageUseCaseDeps
}

func NewSaveChannelMessageUseCase(deps SaveChannelMessageUseCaseDeps) *SaveChannelMessageUseCase {
	return &SaveChannelMessageUseCase{deps: deps}
}

func (u SaveChannelMessageUseCase) Execute(
	ctx context.Context,
	msg *chat_message_entity.ChannelMessage,
) error {
	err := u.deps.Repository.Create(ctx, msg)
	if err != nil {
		return err
	}

	channel, err := u.deps.GetChannelUseCase.Execute(ctx, msg.ChannelId)
	if err != nil {
		return err
	}

	workspaceMember, err := u.deps.GetWorkspaceMemberUseCase.Execute(
		ctx,
		channel.WorkspaceId,
		msg.AuthorId,
	)
	if err != nil {
		return err
	}

	err = u.deps.TimeSeriesRepository.Create(ctx, msg.CreatedAt, entity.MessageSentMetadata{
		WorkspaceId:    channel.WorkspaceId,
		ChannelId:      channel.Id,
		AuthorMemberId: workspaceMember.Id,
	})
	if err != nil {
		return err
	}

	err = u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageKindChannelMessage,
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

	return err
}
