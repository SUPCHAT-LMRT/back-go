package get_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type GetWorkspaceUseCaseDeps struct {
	uberdig.In
	repository.WorkspaceRepository
}
type GetWorkspaceUseCase struct {
	repository repository.WorkspaceRepository
}

func NewGetWorkspaceUseCase(repository GetWorkspaceUseCaseDeps) *GetWorkspaceUseCase {
	return &GetWorkspaceUseCase{repository: repository}
}

func (u GetWorkspaceUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId) (*entity.Workspace, error) {
	return u.repository.GetById(ctx, workspaceId)
}
