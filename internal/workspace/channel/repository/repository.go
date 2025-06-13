package repository

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetById(ctx context.Context, id entity.ChannelId) (*entity.Channel, error)
	List(ctx context.Context, workspaceId workspace_entity.WorkspaceId) ([]*entity.Channel, error)
	CountByWorkspaceId(ctx context.Context, id workspace_entity.WorkspaceId) (uint, error)
	UpdateIndex(ctx context.Context, id entity.ChannelId, index int) error
	Delete(ctx context.Context, id entity.ChannelId) error
	ListPrivateChannelsByUser(
		ctx context.Context,
		workspaceId workspace_entity.WorkspaceId,
		memberId workspace_member_entity.WorkspaceMemberId,
	) ([]*entity.Channel, error)
	ListMembersOfPrivateChannel(
		ctx context.Context,
		channelId entity.ChannelId,
	) ([]workspace_member_entity.WorkspaceMemberId, error)
}
