package list_workpace_members

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type ListWorkspaceMembersUseCase struct {
	repository repository.WorkspaceRepository
}

func NewListWorkspaceMembersUseCase(repository repository.WorkspaceRepository) *ListWorkspaceMembersUseCase {
	return &ListWorkspaceMembersUseCase{repository: repository}
}

func (u ListWorkspaceMembersUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId) ([]*entity.WorkspaceMember, error) {
	return u.repository.ListMembers(ctx, workspaceId)
}
