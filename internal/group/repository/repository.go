package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

var (
	MemberAlreadyInGroupErr = errors.New("member already in group")
	GroupNotFoundErr        = errors.New("group not found")
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group, ownerMember *entity.GroupMember) error
	GetGroup(ctx context.Context, groupId entity.GroupId) (*entity.Group, error)
	ListRecentGroups(ctx context.Context) ([]*entity.Group, error)
	Exists(ctx context.Context, groupId entity.GroupId) (bool, error)
	AddMember(ctx context.Context, groupId entity.GroupId, userId user_entity.UserId) error
	ListMembers(ctx context.Context, groupId entity.GroupId) ([]*entity.GroupMember, error)
	UpdateGroupName(ctx context.Context, groupId entity.GroupId, name string) error
}
