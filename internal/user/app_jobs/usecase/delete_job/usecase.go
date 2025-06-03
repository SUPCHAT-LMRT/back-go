package delete_job

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type DeleteJobUseCase struct {
	repo repository.JobRepository
}

func NewDeleteJobUseCase(repo repository.JobRepository) *DeleteJobUseCase {
	return &DeleteJobUseCase{repo: repo}
}

func (uc *DeleteJobUseCase) Execute(ctx context.Context, jobId entity.JobId) error {
	// Vérifier si le job existe
	job, err := uc.repo.FindById(ctx, jobId)
	if err != nil {
		return fmt.Errorf("error finding job: %w", err)
	}
	if job == nil {
		return fmt.Errorf("job with ID '%s' not found", jobId)
	}

	// Supprimer le job
	err = uc.repo.Delete(ctx, jobId)
	if err != nil {
		return fmt.Errorf("error deleting job: %w", err)
	}

	return nil
}
