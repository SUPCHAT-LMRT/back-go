package repository

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type JobRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Job, error)
	FindById(ctx context.Context, jobId entity.JobId) (*entity.Job, error)
	Create(ctx context.Context, job *entity.Job) error
	Delete(ctx context.Context, jobId entity.JobId) error
	Update(ctx context.Context, job *entity.Job) error
	FindAll(ctx context.Context) ([]*entity.Job, error)
	AssignToUser(ctx context.Context, jobId entity.JobId, userId user_entity.UserId) error
	UnassignFromUser(ctx context.Context, jobId entity.JobId, userId user_entity.UserId) error
	EnsureAdminRoleExists(ctx context.Context) (*entity.Job, error)
	FindByUserId(ctx context.Context, userId string) ([]*entity.Job, error)
}
