package list_private_channels

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetPrivateChannelsUseCase struct {
	repository repository.ChannelRepository
}

func NewGetPrivateChannelsUseCase(
	channelRepository repository.ChannelRepository,
) *GetPrivateChannelsUseCase {
	return &GetPrivateChannelsUseCase{repository: channelRepository}
}

func (u *GetPrivateChannelsUseCase) Execute(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
	userId string,
) ([]*entity.Channel, error) {
	return u.repository.ListPrivateChannelsByUser(ctx, workspaceId, userId)
}
