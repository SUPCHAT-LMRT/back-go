package repository

import (
	"context"
	"errors"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

var ErrWorkspaceNotFound = errors.New("workspace not found")

type WorkspaceRepository interface {
	Create(
		ctx context.Context,
		workspace *entity.Workspace,
		ownerMember *entity2.WorkspaceMember,
	) error
	GetById(ctx context.Context, id entity.WorkspaceId) (*entity.Workspace, error)
	ExistsById(ctx context.Context, id entity.WorkspaceId) (bool, error)
	List(ctx context.Context) ([]*entity.Workspace, error)
	ListPublics(ctx context.Context) ([]*entity.Workspace, error)
	ListByUserId(ctx context.Context, userId user_entity.UserId) ([]*entity.Workspace, error)
	Update(ctx context.Context, workspace *entity.Workspace) error
	Delete(ctx context.Context, id entity.WorkspaceId) error
	GetMemberId(
		ctx context.Context,
		workspaceId entity.WorkspaceId,
		userId user_entity.UserId,
	) (entity2.WorkspaceMemberId, error)
}
