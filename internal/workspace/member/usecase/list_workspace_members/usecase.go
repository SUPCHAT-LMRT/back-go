package list_workspace_members

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

type ListWorkspaceMembersUseCase struct {
	repository repository2.WorkspaceMemberRepository
}

func NewListWorkspaceMembersUseCase(
	repository repository2.WorkspaceMemberRepository,
) *ListWorkspaceMembersUseCase {
	return &ListWorkspaceMembersUseCase{repository: repository}
}

func (u ListWorkspaceMembersUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	limit, page int,
) (totalMembers uint, members []*entity2.WorkspaceMember, err error) {
	return u.repository.ListMembers(ctx, workspaceId, limit, page)
}
