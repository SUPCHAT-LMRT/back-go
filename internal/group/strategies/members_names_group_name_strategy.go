package strategies

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"strings"
)

type MembersNamesDefaultGroupNameStrategy struct {
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewMembersNamesGroupNameStrategy(getUserByIdUseCase *get_by_id.GetUserByIdUseCase) DefaultGroupNameStrategy {
	return &MembersNamesDefaultGroupNameStrategy{getUserByIdUseCase: getUserByIdUseCase}
}

func (s MembersNamesDefaultGroupNameStrategy) Handle(ctx context.Context, group *group_entity.Group, members []*group_entity.GroupMember) (string, error) {
	builder := strings.Builder{}

	for i, member := range members {
		user, err := s.getUserByIdUseCase.Execute(ctx, member.UserId)
		if err != nil {
			return "", err
		}

		builder.WriteString(user.Pseudo)
		if i != len(members)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String(), nil
}
