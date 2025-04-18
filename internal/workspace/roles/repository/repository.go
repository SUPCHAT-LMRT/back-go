package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetById(ctx context.Context, roleId string) (*entity.Role, error)
	GetList(ctx context.Context, workspaceId string) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, roleId string) error
}
