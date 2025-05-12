package get_job_for_user

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type GetJobForUserUseCase struct {
	jobRepo repository.JobRepository
}

func NewGetJobForUserUseCase(jobRepo repository.JobRepository) *GetJobForUserUseCase {
	return &GetJobForUserUseCase{jobRepo: jobRepo}
}

func (u *GetJobForUserUseCase) Execute(ctx context.Context, userId string) ([]*entity.Job, error) {
	// Récupérer tous les jobs
	allJobs, err := u.jobRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Récupérer les jobs assignés à l'utilisateur
	userJobs, err := u.jobRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Créer un map pour identifier les jobs assignés
	assignedJobs := make(map[entity.JobId]bool)
	for _, job := range userJobs {
		assignedJobs[job.Id] = true
	}

	// Mettre à jour le champ IsAssigned pour tous les jobs
	for _, job := range allJobs {
		job.IsAssigned = assignedJobs[job.Id]
	}

	return allJobs, nil
}
