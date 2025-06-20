package is_user_in_workspace

import (
	"context"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

type IsUserInWorkspaceUseCase struct {
	repository repository.WorkspaceMemberRepository
}

func NewIsUserInWorkspaceUseCase(
	workspaceMemberRepository repository.WorkspaceMemberRepository,
) *IsUserInWorkspaceUseCase {
	return &IsUserInWorkspaceUseCase{repository: workspaceMemberRepository}
}

func (u *IsUserInWorkspaceUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	userId user_entity.UserId,
) (bool, error) {
	return u.repository.IsMemberByUserIdExists(ctx, workspaceId, userId)
}
