package get_job_for_user

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"sort"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type GetJobForUserUseCase struct {
	jobRepo repository.JobRepository
}

func NewGetJobForUserUseCase(jobRepo repository.JobRepository) *GetJobForUserUseCase {
	return &GetJobForUserUseCase{jobRepo: jobRepo}
}

func (u *GetJobForUserUseCase) Execute(ctx context.Context, userId user_entity.UserId) ([]*entity.Job, error) {
	allJobs, err := u.jobRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	userJobs, err := u.jobRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	assignedJobs := make(map[entity.JobId]bool)
	for _, job := range userJobs {
		assignedJobs[job.Id] = true
	}

	for _, job := range allJobs {
		job.IsAssigned = assignedJobs[job.Id]
	}

	sort.Slice(allJobs, func(i, j int) bool {
		return allJobs[i].Name < allJobs[j].Name
	})

	return allJobs, nil
}
