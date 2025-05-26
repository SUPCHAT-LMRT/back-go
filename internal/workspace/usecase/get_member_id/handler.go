package get_member_id

import (
	"context"
	"errors"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type GetMemberIdUsecase interface {
	Execute(ctx context.Context, workspaceId string, userId string) (string, error)
}

type getMemberIdUsecase struct {
	Repository repository.WorkspaceRepository
}

func NewGetMemberIdUsecase(repo repository.WorkspaceRepository) GetMemberIdUsecase {
	return &getMemberIdUsecase{Repository: repo}
}

func (u *getMemberIdUsecase) Execute(
	ctx context.Context,
	workspaceId string,
	userId string,
) (string, error) {
	workspaceIdEntity := entity.WorkspaceId(workspaceId)
	userIdEntity := user_entity.UserId(userId)

	memberId, err := u.Repository.GetMemberId(ctx, workspaceIdEntity, userIdEntity)
	if err != nil {
		return "", err
	}
	if memberId == "" {
		return "", errors.New("member not found")
	}
	return string(memberId), nil
}
