package get_polls_list

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetPollsListUseCase struct {
	repo repository.PollRepository
}

func NewGetPollsListUseCase(repo repository.PollRepository) *GetPollsListUseCase {
	return &GetPollsListUseCase{repo: repo}
}

func (uc *GetPollsListUseCase) Execute(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
) ([]*entity.Poll, error) {
	fmt.Println("GetPollsListUseCase.Execute called with workspaceId:", workspaceId)
	polls, err := uc.repo.GetAllByWorkspace(ctx, workspaceId)
	if err != nil {
		return nil, err
	}
	return polls, nil
}
