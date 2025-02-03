package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetById(ctx context.Context, id entity.ChannelId) (*entity.Channel, error)
	List(ctx context.Context, workspaceId workspace_entity.WorkspaceId) ([]*entity.Channel, error)
}
