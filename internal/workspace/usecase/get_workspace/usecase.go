package get_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type GetWorkspaceUseCase struct {
	repository repository.WorkspaceRepository
}

func NewGetWorkspaceUseCase(repository repository.WorkspaceRepository) *GetWorkspaceUseCase {
	return &GetWorkspaceUseCase{repository: repository}
}

func (u GetWorkspaceUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId) (*entity.Workspace, error) {
	return u.repository.GetById(ctx, workspaceId)
}
