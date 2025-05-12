package unassign_job

import (
	"context"
	"fmt"
	job_entity "github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type UnassignJobUseCase struct {
	repo repository.JobRepository
}

func NewUnassignJobUseCase(repo repository.JobRepository) *UnassignJobUseCase {
	return &UnassignJobUseCase{repo: repo}
}

func (uc *UnassignJobUseCase) Execute(ctx context.Context, jobId job_entity.JobId, userId entity.UserId) error {
	// Vérifier si le job existe
	job, err := uc.repo.FindById(ctx, jobId)
	if err != nil {
		return fmt.Errorf("error finding job: %w", err)
	}
	if job == nil {
		return fmt.Errorf("job with ID '%s' not found", jobId)
	}
	// Désassigner le job
	err = uc.repo.UnassignFromUser(ctx, jobId, userId)
	if err != nil {
		return fmt.Errorf("error unassigning job: %w", err)
	}

	return nil
}
