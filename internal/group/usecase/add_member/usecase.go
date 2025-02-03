package add_member

import (
	"context"
	"errors"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/group/strategies"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

var (
	GroupNotFoundErr = errors.New("group not found")
)

type AddMemberToGroupUseCaseDeps struct {
	uberdig.In
	Repository               repository.GroupRepository
	DefaultGroupNameStrategy strategies.DefaultGroupNameStrategy
}

type AddMemberToGroupUseCase struct {
	deps AddMemberToGroupUseCaseDeps
}

func NewAddMemberToGroupUseCase(deps AddMemberToGroupUseCaseDeps) *AddMemberToGroupUseCase {
	return &AddMemberToGroupUseCase{deps: deps}
}

// Todo transaction to be sure that if a problem occurs when adding the member to the group, the group is not created
func (u *AddMemberToGroupUseCase) Execute(ctx context.Context, groupId *group_entity.GroupId, inviterUserId, inviteeUserId entity.UserId) (*group_entity.Group, error) {
	var group *group_entity.Group
	var groupMembers []*group_entity.GroupMember
	if groupId == nil || *groupId == "" {
		// First create the group and then add the member
		group = &group_entity.Group{OwnerUserId: inviterUserId}
		groupMembers = []*group_entity.GroupMember{{UserId: inviterUserId}}
		defaultGroupName, err := u.deps.DefaultGroupNameStrategy.Handle(ctx, group, groupMembers)
		if err != nil {
			return nil, err
		}

		group.Name = defaultGroupName

		err = u.deps.Repository.Create(ctx, group, &group_entity.GroupMember{UserId: inviterUserId})
		if err != nil {
			return nil, err
		}

		groupId = &group.Id
	} else {
		// If the group is just created, it is not necessary to check if it exists
		var err error
		group, err = u.deps.Repository.GetGroup(ctx, *groupId)
		if err != nil {
			return nil, err
		}

		groupMembers, err = u.deps.Repository.ListMembers(ctx, *groupId)
		if err != nil {
			return nil, err
		}
	}

	// Add the member to the groupMembers
	groupMembers = append(groupMembers, &group_entity.GroupMember{UserId: inviteeUserId})

	err := u.deps.Repository.AddMember(ctx, *groupId, inviteeUserId)
	if err != nil {
		return nil, err
	}

	defaultGroupName, err := u.deps.DefaultGroupNameStrategy.Handle(ctx, group, groupMembers)
	if err != nil {
		return nil, err
	}

	group.Name = defaultGroupName

	err = u.deps.Repository.UpdateGroupName(ctx, *groupId, defaultGroupName)
	if err != nil {
		return nil, err
	}

	return group, nil
}
