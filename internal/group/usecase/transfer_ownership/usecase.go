package transfer_ownership

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	uberdig "go.uber.org/dig"
)

type TransferGroupOwnershipUseCaseDeps struct {
	uberdig.In
	GroupRepository repository.GroupRepository
	Observers       []TransferGroupOwnershipObserver `group:"transfer_group_ownership_observers"`
}

type TransferGroupOwnershipUseCase struct {
	deps TransferGroupOwnershipUseCaseDeps
}

func NewTransferGroupOwnershipUseCase(deps TransferGroupOwnershipUseCaseDeps) *TransferGroupOwnershipUseCase {
	return &TransferGroupOwnershipUseCase{
		deps: deps,
	}
}

func (uc *TransferGroupOwnershipUseCase) Execute(ctx context.Context, groupId group_entity.GroupId, newOwnerId group_entity.GroupMemberId) error {
	// Fetch the group to ensure it exists
	group, err := uc.deps.GroupRepository.GetGroup(ctx, groupId)
	if err != nil {
		return err
	}

	// Check if the new owner is a member of the group
	isMember, err := uc.deps.GroupRepository.IsMember(ctx, groupId, newOwnerId)
	if err != nil {
		return err
	}
	if !isMember {
		return repository.ErrMemberNotFound
	}

	// Update the group's owner
	err = uc.deps.GroupRepository.TransferOwnership(ctx, groupId, newOwnerId)
	if err != nil {
		return err
	}

	// Notify observers about the ownership transfer
	for _, observer := range uc.deps.Observers {
		observer.NotifyOwnershipTransferred(group, newOwnerId)
	}

	return nil
}
