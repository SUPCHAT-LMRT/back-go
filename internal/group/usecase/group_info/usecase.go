package group_info

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	list_members "github.com/supchat-lmrt/back-go/internal/group/usecase/list_members_users"
	uberdig "go.uber.org/dig"
	"time"
)

type GetGroupInfoUseCaseDeps struct {
	uberdig.In
	GroupRepository         repository.GroupRepository
	ListGroupMembersUseCase *list_members.ListGroupMembersUseCase
}

type GetGroupInfoUseCase struct {
	deps GetGroupInfoUseCaseDeps
}

func NewGetGroupInfoUseCase(deps GetGroupInfoUseCaseDeps) *GetGroupInfoUseCase {
	return &GetGroupInfoUseCase{
		deps: deps,
	}
}

func (uc *GetGroupInfoUseCase) Execute(ctx context.Context, groupId group_entity.GroupId) (*GroupInfo, error) {
	group, err := uc.deps.GroupRepository.GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	// Fetch group members
	members, err := uc.deps.ListGroupMembersUseCase.Execute(ctx, groupId)
	if err != nil {
		return nil, err
	}

	// Create and return group info
	groupInfo := &GroupInfo{
		Id:        group.Id,
		Name:      group.Name,
		Members:   members,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

	return groupInfo, nil
}

type GroupInfo struct {
	Id        group_entity.GroupId
	Name      string
	Members   []*list_members.ListGroupMembersResponse
	CreatedAt time.Time
	UpdatedAt time.Time
}
