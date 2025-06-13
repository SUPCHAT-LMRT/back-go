package repository

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

var (
	ErrMemberAlreadyInGroup = errors.New("member already in group")
	ErrGroupNotFound        = errors.New("group not found")
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group, ownerUserId user_entity.UserId) error
	GetGroup(ctx context.Context, groupId entity.GroupId) (*entity.Group, error)
	GetMemberByUserId(ctx context.Context, groupId entity.GroupId, userId user_entity.UserId) (*entity.GroupMember, error)
	DeleteGroup(ctx context.Context, groupId entity.GroupId) error
	ListRecentGroups(ctx context.Context, userId user_entity.UserId) ([]*entity.Group, error)
	Exists(ctx context.Context, groupId entity.GroupId) (bool, error)
	AddMember(ctx context.Context, groupId entity.GroupId, userId user_entity.UserId) error
	RemoveMember(ctx context.Context, groupId entity.GroupId, memberId entity.GroupMemberId) error
	ListMembers(ctx context.Context, groupId entity.GroupId) ([]*entity.GroupMember, error)
	UpdateGroupName(ctx context.Context, groupId entity.GroupId, name string) error
	TransferOwnership(ctx context.Context, groupId entity.GroupId, newOwnerMemberId entity.GroupMemberId) error
}
