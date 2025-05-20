package update_banner

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

type UpdateWorkspaceBannerUseCaseDeps struct {
	uberdig.In
	Strategy   UpdateWorkspaceBannerStrategy
	Repository repository.WorkspaceRepository
	Observers  []SaveBannerWorkspaceObserver `group:"save_banner_workspace_observers"`
}

type UpdateWorkspaceBannerUseCase struct {
	deps UpdateWorkspaceBannerUseCaseDeps
}

func NewUpdateWorkspaceBannerUseCase(deps UpdateWorkspaceBannerUseCaseDeps) *UpdateWorkspaceBannerUseCase {
	return &UpdateWorkspaceBannerUseCase{deps: deps}
}

func (u *UpdateWorkspaceBannerUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId, image UpdateImage) error {
	workspace, err := u.deps.Repository.GetById(ctx, workspaceId)
	if err != nil {
		return err
	}

	err = u.deps.Strategy.Handle(ctx, workspaceId, image.ImageReader, image.ContentType)
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
