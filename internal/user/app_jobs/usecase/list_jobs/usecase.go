package list_jobs

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type ListJobsUseCase struct {
	repo repository.JobRepository
}

func NewListJobsUseCase(repo repository.JobRepository) *ListJobsUseCase {
	return &ListJobsUseCase{repo: repo}
}

func (uc *ListJobsUseCase) Execute(ctx context.Context) ([]*entity.Job, error) {
	return uc.repo.FindAll(ctx)
}
