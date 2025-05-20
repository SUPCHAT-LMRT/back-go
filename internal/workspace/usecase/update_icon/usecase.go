package update_icon

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
	"io"
)

var (
	WorkspaceNotFoundErr = repository.WorkspaceNotFoundErr
)

type UpdateWorkspaceIconUseCaseDeps struct {
	uberdig.In
	Strategy   UpdateWorkspaceIconStrategy
	Repository repository.WorkspaceRepository
	Observers  []SaveIconWorkspaceObserver `group:"update_icon_workspace_observers"`
}

type UpdateWorkspaceIconUseCase struct {
	deps UpdateWorkspaceIconUseCaseDeps
}

func NewUpdateWorkspaceIconUseCase(deps UpdateWorkspaceIconUseCaseDeps) *UpdateWorkspaceIconUseCase {
	return &UpdateWorkspaceIconUseCase{deps: deps}
}

func (u *UpdateWorkspaceIconUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, image UpdateImage) error {
	exists, err := u.deps.Repository.ExistsById(ctx, workspaceId)
	if err != nil {
		return err
	}
	if !exists {
		return WorkspaceNotFoundErr
	}

	if err := u.deps.Strategy.Handle(ctx, workspaceId, image.ImageReader, image.ContentType); err != nil {
		return err
	}

	workspace, err := u.deps.Repository.GetById(ctx, workspaceId)
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyUpdateBannerWorkspace(workspace)
	}

	return nil
}

type UpdateImage struct {
	ImageReader io.Reader
	ContentType string
}
