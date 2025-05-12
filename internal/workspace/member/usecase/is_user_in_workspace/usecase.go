package is_user_in_workspace

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

type IsUserInWorkspaceUseCase struct {
	repository repository.WorkspaceMemberRepository
}

func NewIsUserInWorkspaceUseCase(repository repository.WorkspaceMemberRepository) *IsUserInWorkspaceUseCase {
	return &IsUserInWorkspaceUseCase{repository: repository}
}

func (u *IsUserInWorkspaceUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (bool, error) {
	memberId := entity2.WorkspaceMemberId(userId)
	return u.repository.IsMemberExists(ctx, workspaceId, memberId)
}
