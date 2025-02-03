package get_workpace_member

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type GetWorkspaceMemberUseCase struct {
	repository repository.WorkspaceRepository
}

func NewGetWorkspaceMemberUseCase(repository repository.WorkspaceRepository) *GetWorkspaceMemberUseCase {
	return &GetWorkspaceMemberUseCase{repository: repository}
}

func (u GetWorkspaceMemberUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (*entity.WorkspaceMember, error) {
	return u.repository.GetMemberByUserId(ctx, workspaceId, userId)
}
