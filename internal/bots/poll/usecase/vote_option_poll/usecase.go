package vote_option_poll

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
)

// CustomError est une structure d'erreur personnalisée avec un code
type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

const (
	ErrorCodeAlreadyVoted = "ALREADY_VOTED"
	ErrorCodePollNotFound = "POLL_NOT_FOUND"
	ErrorCodeUpdateFailed = "UPDATE_FAILED"
)

type VoteOptionPollUseCase struct {
	repo repository.PollRepository
}

func NewVoteOptionPollUseCase(repo repository.PollRepository) *VoteOptionPollUseCase {
	return &VoteOptionPollUseCase{repo: repo}
}

func (uc *VoteOptionPollUseCase) Execute(ctx context.Context, pollId string, optionId string, userId string) error {
	poll, err := uc.repo.GetById(ctx, pollId)
	if err != nil {
		return &CustomError{Code: ErrorCodePollNotFound, Message: "Sondage non trouvé"}
	}

	for _, opt := range poll.Options {
		for _, voter := range opt.Voters {
			if voter == userId {
				return &CustomError{Code: ErrorCodeAlreadyVoted, Message: "Vous avez déjà voté pour ce sondage"}
			}
		}
	}

	var optionFound bool
	for i, opt := range poll.Options {
		if opt.Id == optionId {
			optionFound = true
			poll.Options[i].Voters = append(poll.Options[i].Voters, userId)
			poll.Options[i].Votes++
			break
		}
	}
	if !optionFound {
		return &CustomError{Code: ErrorCodeUpdateFailed, Message: "Option non trouvée"}
	}

	err = uc.repo.Vote(ctx, poll)
	if err != nil {
		return &CustomError{Code: ErrorCodeUpdateFailed, Message: "Erreur lors de l'enregistrement du vote"}
	}

	return nil
}
