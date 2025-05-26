package get_workspace

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type GetWorkspaceUseCaseDeps struct {
	uberdig.In
	repository repository.WorkspaceRepository
}
type GetWorkspaceUseCase struct {
	deps GetWorkspaceUseCaseDeps
}

func NewGetWorkspaceUseCase(deps GetWorkspaceUseCaseDeps) *GetWorkspaceUseCase {
	return &GetWorkspaceUseCase{deps: deps}
}

func (u GetWorkspaceUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
) (*entity.Workspace, error) {
	return u.deps.repository.GetById(ctx, workspaceId)
}
