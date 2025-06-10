package repository

import (
	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
)

type MongoPollMapper struct{}

func NewMongoPollMapper() *MongoPollMapper {
	return &MongoPollMapper{}
}

func (m *MongoPollMapper) MapFromEntity(poll *entity.Poll) (*MongoPoll, error) {
	options := make([]Option, len(poll.Options))
	for i, opt := range poll.Options {
		voters := opt.Voters
		if voters == nil {
			voters = []string{}
		}
		options[i] = Option{
			Id:     opt.Id,
			Text:   opt.Text,
			Votes:  opt.Votes,
			Voters: voters,
		}
	}

	return &MongoPoll{
		Id:          poll.Id,
		Question:    poll.Question,
		Options:     options,
		CreatedBy:   poll.CreatedBy,
		WorkspaceId: poll.WorkspaceId,
		CreatedAt:   poll.CreatedAt,
		ExpiresAt:   poll.ExpiresAt,
	}, nil
}

func (m *MongoPollMapper) MapToEntity(mongoPoll *MongoPoll) (*entity.Poll, error) {
	options := make([]entity.Option, len(mongoPoll.Options))
	for i, opt := range mongoPoll.Options {
		options[i] = entity.Option{
			Id:     opt.Id,
			Text:   opt.Text,
			Votes:  opt.Votes,
			Voters: opt.Voters,
		}
	}

	return &entity.Poll{
		Id:          mongoPoll.Id,
		Question:    mongoPoll.Question,
		Options:     options,
		CreatedBy:   mongoPoll.CreatedBy,
		WorkspaceId: mongoPoll.WorkspaceId,
		CreatedAt:   mongoPoll.CreatedAt,
		ExpiresAt:   mongoPoll.ExpiresAt,
	}, nil
}
