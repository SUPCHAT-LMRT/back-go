package permissions

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
)

type CheckPermissionJobUseCase struct {
	jobRepository repository.JobRepository
}

func NewCheckPermissionJobUseCase(
	jobRepository repository.JobRepository,
) *CheckPermissionJobUseCase {
	return &CheckPermissionJobUseCase{jobRepository: jobRepository}
}

func (u *CheckPermissionJobUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
	permission uint64,
) (bool, error) {
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
