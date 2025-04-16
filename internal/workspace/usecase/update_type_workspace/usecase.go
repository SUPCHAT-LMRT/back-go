package update_type_workspace

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type UpdateTypeWorkspaceUseCaseDeps struct {
	uberdig.In
	WorkspaceRepository repository.WorkspaceRepository
	Observers           []SaveTypeWorkspaceObserver `group:"update_type_workspace_observers"`
}

type UpdateTypeWorkspaceUseCase struct {
	deps UpdateTypeWorkspaceUseCaseDeps
}

func NewUpdateTypeWorkspaceUseCase(deps UpdateTypeWorkspaceUseCaseDeps) *UpdateTypeWorkspaceUseCase {
	return &UpdateTypeWorkspaceUseCase{deps: deps}
}

func (u *UpdateTypeWorkspaceUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, typeWorkspace entity.WorkspaceType) error {
	if typeWorkspace == "" {
		return errors.New("type_workspace is required")
	}

	workspace, err := u.deps.WorkspaceRepository.GetById(ctx, workspaceId)
	if err != nil {
		return err
	}

	workspace.Type = typeWorkspace

	if err := u.deps.WorkspaceRepository.Update(ctx, workspace); err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyUpdateTypeWorkspace(workspace)
	}

	return nil

}
