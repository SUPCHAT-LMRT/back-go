package unvote_option_poll

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
)

type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

const (
	ErrorCodeNotVoted     = "NOT_VOTED"
	ErrorCodePollNotFound = "POLL_NOT_FOUND"
	ErrorCodeUpdateFailed = "UPDATE_FAILED"
)

type UnvoteOptionPollUseCase struct {
	repo repository.PollRepository
}

func NewUnvoteOptionPollUseCase(repo repository.PollRepository) *UnvoteOptionPollUseCase {
	return &UnvoteOptionPollUseCase{repo: repo}
}

func (uc *UnvoteOptionPollUseCase) Execute(
	ctx context.Context,
	pollId string,
	userId string,
) error {
	poll, err := uc.repo.GetById(ctx, pollId)
	if err != nil {
		return &CustomError{Code: ErrorCodePollNotFound, Message: "Sondage non trouvé"}
	}

	var found bool
	for i, opt := range poll.Options {
		for j, voter := range opt.Voters {
			if voter == userId {
				// Retirer l'utilisateur de la liste des votants
				poll.Options[i].Voters = append(opt.Voters[:j], opt.Voters[j+1:]...)
				if poll.Options[i].Votes > 0 {
					poll.Options[i].Votes--
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return &CustomError{
			Code:    ErrorCodeNotVoted,
			Message: "Vous n'avez pas voté pour ce sondage",
		}
	}

	err = uc.repo.Vote(ctx, poll)
	if err != nil {
		return &CustomError{
			Code:    ErrorCodeUpdateFailed,
			Message: "Erreur lors de la suppression du vote",
		}
	}
	return nil
}
