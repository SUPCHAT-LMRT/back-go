package update_banner

import (
	"context"
	"io"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

var WorkspaceNotFoundErr = repository.ErrWorkspaceNotFound

type UpdateWorkspaceBannerUseCaseDeps struct {
	uberdig.In
	Strategy   UpdateWorkspaceBannerStrategy
	Repository repository.WorkspaceRepository
}

type UpdateWorkspaceBannerUseCase struct {
	deps UpdateWorkspaceBannerUseCaseDeps
}

func NewUpdateWorkspaceBannerUseCase(
	deps UpdateWorkspaceBannerUseCaseDeps,
) *UpdateWorkspaceBannerUseCase {
	return &UpdateWorkspaceBannerUseCase{deps: deps}
}

func (u *UpdateWorkspaceBannerUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	image UpdateImage,
) error {
	exists, err := u.deps.Repository.ExistsById(ctx, workspaceId)
	if err != nil {
		return err
	}
	if !exists {
		return WorkspaceNotFoundErr
	}

	return u.deps.Strategy.Handle(ctx, workspaceId, image.ImageReader, image.ContentType)
}

type UpdateImage struct {
	ImageReader io.Reader
	ContentType string
}
