package create_group

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type CreateGroupUseCase struct {
	groupRepository repository.GroupRepository
}

func NewCreateGroupUseCase(groupRepository repository.GroupRepository) *CreateGroupUseCase {
	return &CreateGroupUseCase{groupRepository: groupRepository}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, input CreateGroupInput) (*group_entity.Group, error) {
	now := time.Now()
	group := group_entity.Group{
		Name:        input.GroupName,
		OwnerUserId: input.OwnerUserId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := uc.groupRepository.Create(ctx, &group)
	if err != nil {
		return nil, err
	}

	// TODO impl meilisearch
	// TODO impl sync recent chats

	return &group, nil
}

type CreateGroupInput struct {
	OwnerUserId user_entity.UserId
	GroupName   string
	UsersIds    []user_entity.UserId
}
