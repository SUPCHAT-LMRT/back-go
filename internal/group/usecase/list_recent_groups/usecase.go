package list_recent_groups

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	uberdig "go.uber.org/dig"
)

type ListRecentGroupsUseCaseDeps struct {
	uberdig.In
	Repository repository.GroupRepository
}

type ListRecentGroupsUseCase struct {
	deps ListRecentGroupsUseCaseDeps
}

func NewListRecentGroupsUseCase(deps ListRecentGroupsUseCaseDeps) *ListRecentGroupsUseCase {
	return &ListRecentGroupsUseCase{deps: deps}
}

func (u *ListRecentGroupsUseCase) Execute(ctx context.Context) ([]*entity.Group, error) {
	return u.deps.Repository.ListRecentGroups(ctx)
}
