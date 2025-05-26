package create_job

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type CreateJobUseCase struct {
	repo repository.JobRepository
}

func NewCreateJobUseCase(repo repository.JobRepository) *CreateJobUseCase {
	return &CreateJobUseCase{repo: repo}
}

func (uc *CreateJobUseCase) Execute(ctx context.Context, name string) (*entity.Job, error) {
	// Vérifier si un job avec le même nom existe déjà
	existingJob, err := uc.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existingJob != nil {
		return nil, fmt.Errorf("a job with the name '%s' already exists", name)
	}

	job := &entity.Job{
		Name:        name,
		Permissions: 0, // Pas de permissions fonctionnelles
	}

	err = uc.repo.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}
