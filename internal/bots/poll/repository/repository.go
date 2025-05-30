package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type PollRepository interface {
	Create(ctx context.Context, poll *entity.Poll) error
	GetById(ctx context.Context, pollId string) (*entity.Poll, error)
	GetAllByWorkspace(ctx context.Context, workspaceId entity2.WorkspaceId) ([]*entity.Poll, error)
	Delete(ctx context.Context, pollId string) error
	Vote(ctx context.Context, poll *entity.Poll) error
	IncrementVote(ctx context.Context, pollId string, optionId string) error
}
