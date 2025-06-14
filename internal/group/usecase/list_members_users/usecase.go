package list_members

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type ListGroupMembersUseCaseDeps struct {
	uberdig.In
	GroupRepository        repository.GroupRepository
	GetUserByIdUseCase     *get_by_id.GetUserByIdUseCase
	GetPublicStatusUseCase *get_public_status.GetPublicStatusUseCase
}

type ListGroupMembersUseCase struct {
	deps ListGroupMembersUseCaseDeps
}

func NewListGroupMembersUseCase(deps ListGroupMembersUseCaseDeps) *ListGroupMembersUseCase {
	return &ListGroupMembersUseCase{
		deps: deps,
	}
}

func (uc *ListGroupMembersUseCase) Execute(ctx context.Context, groupId group_entity.GroupId) ([]*ListGroupMembersResponse, error) {
	group, err := uc.deps.GroupRepository.GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	members, err := uc.deps.GroupRepository.ListMembers(ctx, groupId)
	if err != nil {
		return nil, err
	}

	users := make([]*ListGroupMembersResponse, len(members))
	for i, member := range members {
		user, err := uc.deps.GetUserByIdUseCase.Execute(ctx, member.UserId)
		if err != nil {
			return nil, err
		}

		status, err := uc.deps.GetPublicStatusUseCase.Execute(ctx, user.Id, entity.StatusOffline)
		if err != nil {
			return nil, err
		}

		users[i] = &ListGroupMembersResponse{
			Id:           member.Id,
			UserId:       user.Id,
			UserName:     user.FullName(),
			Email:        user.Email,
			IsGroupOwner: group.OwnerMemberId == member.Id,
			Status:       status,
		}

	}

	return users, nil
}

type ListGroupMembersResponse struct {
	Id           group_entity.GroupMemberId
	UserId       user_entity.UserId
	UserName     string
	Email        string
	IsGroupOwner bool
	Status       entity.Status
}
