package assign_job

import (
	"context"
	"fmt"
	job_entity "github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type AssignJobUseCase struct {
	repo repository.JobRepository
}

func NewAssignJobUseCase(repo repository.JobRepository) *AssignJobUseCase {
	return &AssignJobUseCase{repo: repo}
}

func (uc *AssignJobUseCase) Execute(ctx context.Context, jobId job_entity.JobId, userId entity.UserId) error {
	// Vérifier si le job existe
	job, err := uc.repo.FindById(ctx, jobId)
	if err != nil {
		return fmt.Errorf("error finding job: %w", err)
	}
	if job == nil {
		return fmt.Errorf("job with ID '%s' not found", jobId)
	}

	// Assigner le job
	err = uc.repo.AssignToUser(ctx, jobId, userId)
	if err != nil {
		return fmt.Errorf("error assigning job: %w", err)
	}

	return nil
}
