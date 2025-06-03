package update_job

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type UpdateJobUseCase struct {
	repo repository.JobRepository
}

func NewUpdateJobUseCase(repo repository.JobRepository) *UpdateJobUseCase {
	return &UpdateJobUseCase{repo: repo}
}

func (uc *UpdateJobUseCase) Execute(
	ctx context.Context,
	jobId entity.JobId,
	name string,
) (*entity.Job, error) {
	// Vérifier si le job existe
	job, err := uc.repo.FindById(ctx, jobId)
	if err != nil {
		return nil, fmt.Errorf("error finding job: %w", err)
	}
	if job == nil {
		return nil, fmt.Errorf("job with ID '%s' not found", jobId)
	}

	// Mettre à jour uniquement le champ Name
	job.Name = name

	// Sauvegarder les modifications
	err = uc.repo.Update(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("error updating job: %w", err)
	}

	return job, nil
}
