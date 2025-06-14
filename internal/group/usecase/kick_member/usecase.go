package kick_member

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/transfer_ownership"
	"github.com/supchat-lmrt/back-go/internal/search/group"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"math/rand/v2"
)

type KickMemberUseCaseDeps struct {
	uberdig.In
	GroupRepository               repository.GroupRepository
	TransferGroupOwnershipUseCase *transfer_ownership.TransferGroupOwnershipUseCase
	SearchGroupSyncManager        group.SearchGroupSyncManager
	Observers                     []KickGroupMemberObserver `group:"kick_group_member_observer"`
}

type KickMemberUseCase struct {
	deps KickMemberUseCaseDeps
}

func NewKickMemberUseCase(deps KickMemberUseCaseDeps) *KickMemberUseCase {
	return &KickMemberUseCase{
		deps: deps,
	}
}

func (uc *KickMemberUseCase) Execute(ctx context.Context, memberIdLeft group_entity.GroupMemberId, groupId group_entity.GroupId) error {
	groupResult, err := uc.deps.GroupRepository.GetGroup(ctx, groupId)
	if err != nil {
		return err
	}

	groupMembers, err := uc.deps.GroupRepository.ListMembers(ctx, groupId)
	if err != nil {
		return err
	}

	var userLeftId user_entity.UserId
	for _, member := range groupMembers {
		if member.Id == memberIdLeft {
			userLeftId = member.UserId
			break
		}
	}

	// If the group has only one member (the owner), we can remove the group entirely.
	if len(groupMembers) == 1 {
		err = uc.deps.GroupRepository.DeleteGroup(ctx, groupId)
		if err != nil {
			return err
		}

		err = uc.deps.SearchGroupSyncManager.RemoveGroup(ctx, groupResult.Id)
		if err != nil {
			return err
		}

		for _, observer := range uc.deps.Observers {
			observer.NotifyGroupMemberKicked(groupResult, memberIdLeft, userLeftId)
		}

		return nil
	}

	// If the member leaving is the owner, we need to handle ownership transfer.
	if groupResult.OwnerMemberId == memberIdLeft {
		// Filter out the member that is leaving
		var potentialOwners []*group_entity.GroupMember
		for _, member := range groupMembers {
			if member.Id != memberIdLeft {
				potentialOwners = append(potentialOwners, member)
			}
		}

		// Find a new owner (shuffle the members and take the first one)
		rand.Shuffle(len(potentialOwners), func(i, j int) {
			potentialOwners[i], potentialOwners[j] = potentialOwners[j], potentialOwners[i]
		})

		err = uc.deps.TransferGroupOwnershipUseCase.Execute(ctx, groupId, potentialOwners[0].Id)
		if err != nil {
			return err
		}
	}

	// We just remove the member.
	err = uc.deps.GroupRepository.RemoveMember(ctx, groupId, memberIdLeft)
	if err != nil {
		return err
	}

	for _, observer := range uc.deps.Observers {
		observer.NotifyGroupMemberKicked(groupResult, memberIdLeft, userLeftId)
	}

	return nil
}
