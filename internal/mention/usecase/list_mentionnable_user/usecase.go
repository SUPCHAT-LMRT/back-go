package list_mentionnable_user

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type (
	Channel struct {
		Id          string
		IsPrivate   bool
		WorkspaceId string
	}
	MentionnableUser struct {
		Id       entity.UserId
		Username string
	}

	ListMentionnableUserUseCase struct {
		channelRepo   repository.ChannelRepository
		workspaceRepo repository2.WorkspaceRepository
	}
)

func NewListMentionnableUserUseCase(
	channelRepo repository.ChannelRepository,
	workspaceRepo repository2.WorkspaceRepository,
) *ListMentionnableUserUseCase {
	return &ListMentionnableUserUseCase{
		channelRepo:   channelRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (u *ListMentionnableUserUseCase) Execute(ctx context.Context, channelId channel_entity.ChannelId) ([]MentionnableUser, error) {
	channel, err := u.channelRepo.GetById(ctx, channelId)
	if err != nil {
		return nil, err
	}

	if channel.IsPrivate {
		return u.channelRepo.GetMembers(ctx, channelId)
	}

	return u.workspaceRepo.GetMembers(ctx, channel.WorkspaceId)
}
