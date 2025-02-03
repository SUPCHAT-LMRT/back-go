package list_channels

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ListChannelsUseCase struct {
	repository repository.ChannelRepository
}

func NewListChannelsUseCase(repository repository.ChannelRepository) *ListChannelsUseCase {
	return &ListChannelsUseCase{repository: repository}
}

func (u *ListChannelsUseCase) Execute(ctx context.Context, workspaceId workspace_entity.WorkspaceId) ([]*entity.Channel, error) {
	return u.repository.List(ctx, workspaceId)
}
