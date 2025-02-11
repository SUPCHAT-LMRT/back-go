package is_user_in_workspace

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type IsUserInWorkspaceUseCase struct {
	repository repository.WorkspaceRepository
}

func NewIsUserInWorkspaceUseCase(repository repository.WorkspaceRepository) *IsUserInWorkspaceUseCase {
	return &IsUserInWorkspaceUseCase{repository: repository}
}

func (u *IsUserInWorkspaceUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (bool, error) {
	return u.repository.IsMemberExists(ctx, workspaceId, userId)
}
