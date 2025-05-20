package permissions

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type CheckPermissionJobUseCase struct {
	jobRepository repository.JobRepository
}

func NewCheckPermissionJobUseCase(jobRepository repository.JobRepository) *CheckPermissionJobUseCase {
	return &CheckPermissionJobUseCase{jobRepository: jobRepository}
}

func (u *CheckPermissionJobUseCase) Execute(ctx context.Context, userId string, permission uint64) (bool, error) {
	jobs, err := u.jobRepository.FindByUserId(ctx, userId)
	if err != nil {
		return false, err
	}

	for _, job := range jobs {
		if job.IsAssigned && job.HasPermission(permission) {
			return true, nil
		}
	}
	return false, nil
}
