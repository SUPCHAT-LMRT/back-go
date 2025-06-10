package get_poll_by_id

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
)

type GetPollByIdUseCase struct {
	repo repository.PollRepository
}

func NewGetPollByIdUseCase(repo repository.PollRepository) *GetPollByIdUseCase {
	return &GetPollByIdUseCase{repo: repo}
}

func (uc *GetPollByIdUseCase) Execute(ctx context.Context, pollId string) (*entity.Poll, error) {
	poll, err := uc.repo.GetById(ctx, pollId)
	if err != nil {
		return nil, errors.New("poll not found")
	}
	return poll, nil
}
