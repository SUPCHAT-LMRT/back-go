package delete_poll

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
)

type DeletePollUseCase struct {
	repo repository.PollRepository
}

func NewDeletePollUseCase(repo repository.PollRepository) *DeletePollUseCase {
	return &DeletePollUseCase{repo: repo}
}

func (uc *DeletePollUseCase) Execute(ctx context.Context, pollId string) error {
	return uc.repo.Delete(ctx, pollId)
}
