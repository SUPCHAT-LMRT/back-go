package create_poll

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type CreatePollUseCase struct {
	repo repository.PollRepository
}

func NewCreatePollUseCase(repo repository.PollRepository) *CreatePollUseCase {
	return &CreatePollUseCase{repo: repo}
}

func (uc *CreatePollUseCase) Execute(ctx context.Context, question string, options []string, createdBy string, workspaceId string, expiresAt time.Time) (*entity.Poll, error) {
	if len(options) < 2 {
		return nil, errors.New("a poll must have at least two options")
	}

	pollOptions := make([]entity.Option, len(options))
	for i, opt := range options {
		pollOptions[i] = entity.Option{
			Id:     bson.NewObjectID().Hex(),
			Text:   opt,
			Votes:  0,
			Voters: []string{},
		}
	}

	poll := &entity.Poll{
		Id:          bson.NewObjectID().Hex(),
		Question:    question,
		Options:     pollOptions,
		CreatedBy:   createdBy,
		WorkspaceId: entity2.WorkspaceId(workspaceId),
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	err := uc.repo.Create(ctx, poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}
