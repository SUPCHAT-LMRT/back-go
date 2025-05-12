package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type JobRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Job, error)
	FindById(ctx context.Context, jobId string) (*entity.Job, error)
	Create(ctx context.Context, job *entity.Job) error
	Delete(ctx context.Context, jobId string) error
	Update(ctx context.Context, job *entity.Job) error
	FindAll(ctx context.Context) ([]*entity.Job, error)
	AssignToUser(ctx context.Context, jobId string, userId user_entity.UserId) error
	UnassignFromUser(ctx context.Context, jobId string, userId user_entity.UserId) error
	EnsureAdminRoleExists(ctx context.Context) error
	FindByUserId(ctx context.Context, userId string) ([]*entity.Job, error)
}
