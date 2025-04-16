package update_info_workspaces

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type UpdateInfoWorkspacesUseCaseDeps struct {
	uberdig.In
	WorkspaceRepository repository.WorkspaceRepository
	Observers           []SaveInfoWorkspacesObserver `group:"update_info_workspaces_observers"`
}

type UpdateInfoWorkspacesUseCase struct {
	deps UpdateInfoWorkspacesUseCaseDeps
}

func NewUpdateInfoWorkspacesUseCase(deps UpdateInfoWorkspacesUseCaseDeps) *UpdateInfoWorkspacesUseCase {
	return &UpdateInfoWorkspacesUseCase{deps: deps}
}

func (u *UpdateInfoWorkspacesUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, name, topic string) error {
	if name == "" {
		return errors.New("name is required")
	}

	workspace, err := u.deps.WorkspaceRepository.GetById(ctx, workspaceId)
	if err != nil {
		return err
	}

	workspace.Name = name
	workspace.Topic = topic

	if err := u.deps.WorkspaceRepository.Update(ctx, workspace); err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyUpdateInfoWorkspaces(workspace)
	}

	return nil
}
