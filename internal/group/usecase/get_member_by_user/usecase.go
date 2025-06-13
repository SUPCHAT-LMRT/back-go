package get_member_by_user

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type GetMemberByUserUseCase struct {
	repository repository.GroupRepository
}

func NewGetMemberByUserUseCase(repository repository.GroupRepository) *GetMemberByUserUseCase {
	return &GetMemberByUserUseCase{
		repository: repository,
	}
}

func (uc *GetMemberByUserUseCase) Execute(ctx context.Context, groupId group_entity.GroupId, userId user_entity.UserId) (*group_entity.GroupMember, error) {
	return uc.repository.GetMemberByUserId(ctx, groupId, userId)
}
