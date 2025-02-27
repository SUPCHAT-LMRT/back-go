package get_workpace_member

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

type GetWorkspaceMemberUseCase struct {
	repository repository2.WorkspaceMemberRepository
}

func NewGetWorkspaceMemberUseCase(repository repository2.WorkspaceMemberRepository) *GetWorkspaceMemberUseCase {
	return &GetWorkspaceMemberUseCase{repository: repository}
}

func (u GetWorkspaceMemberUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (*entity2.WorkspaceMember, error) {
	return u.repository.GetMemberByUserId(ctx, workspaceId, userId)
}
