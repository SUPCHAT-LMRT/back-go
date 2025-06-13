package create_group

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"time"
)

type CreateGroupUseCaseDeps struct {
	uberdig.In
	GroupRepository repository.GroupRepository
	Observers       []GroupCreatedObserver `group:"group_created_observer"`
}

type CreateGroupUseCase struct {
	deps CreateGroupUseCaseDeps
}

func NewCreateGroupUseCase(deps CreateGroupUseCaseDeps) *CreateGroupUseCase {
	return &CreateGroupUseCase{deps: deps}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, input CreateGroupInput) (*group_entity.Group, error) {
	now := time.Now()
	group := group_entity.Group{
		Name:      input.GroupName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := uc.deps.GroupRepository.Create(ctx, &group, input.OwnerUserId)
	if err != nil {
		return nil, err
	}

	for _, id := range input.UsersIds {
		if id == input.OwnerUserId {
			// Skip adding the owner as a member, they are already the owner
			continue
		}

		err = uc.deps.GroupRepository.AddMember(ctx, group.Id, id)
		if err != nil {
			return nil, err
		}
	}

	// TODO impl meilisearch
	for _, observer := range uc.deps.Observers {
		observer.NotifyGroupMemberAdded(&group)
	}

	return &group, nil
}

type CreateGroupInput struct {
	OwnerUserId user_entity.UserId
	GroupName   string
	UsersIds    []user_entity.UserId
}
