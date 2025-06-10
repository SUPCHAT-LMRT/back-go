package count_channels

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type CountChannelsUseCase struct {
	repository repository.ChannelRepository
}

func NewCountChannelsUseCase(channelRepository repository.ChannelRepository) *CountChannelsUseCase {
	return &CountChannelsUseCase{repository: channelRepository}
}

func (u CountChannelsUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
) (uint, error) {
	return u.repository.CountByWorkspaceId(ctx, workspaceId)
}
