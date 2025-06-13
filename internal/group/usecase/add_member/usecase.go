package add_member

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type AddMemberToGroupUseCaseDeps struct {
	uberdig.In
	Repository repository.GroupRepository
	Observers  []AddGroupMemberObserver `group:"add_group_member_observer"`
}

type AddMemberToGroupUseCase struct {
	deps AddMemberToGroupUseCaseDeps
}

func NewAddMemberToGroupUseCase(deps AddMemberToGroupUseCaseDeps) *AddMemberToGroupUseCase {
	return &AddMemberToGroupUseCase{deps: deps}
}

func (u *AddMemberToGroupUseCase) Execute(
	ctx context.Context,
	givenGroupId group_entity.GroupId,
	inviteeUserId entity.UserId,
) error {
	groupId := givenGroupId

	// Check if the group exists
	group, err := u.deps.Repository.GetGroup(ctx, groupId)
	if err != nil {
		return err
	}

	err = u.deps.Repository.AddMember(ctx, groupId, inviteeUserId)
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyGroupMemberKicked(group, inviteeUserId)
	}

	return nil
}
