package add_reaction

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
)

type AddReactionUseCase struct {
	repository repository.ChannelMessageRepository
}

func NewAddReactionUseCase(repository repository.ChannelMessageRepository) *AddReactionUseCase {
	return &AddReactionUseCase{repository: repository}
}

func (u *AddReactionUseCase) Execute(ctx context.Context, reaction entity.ChannelMessageReaction) error {
	return u.repository.AddReaction(ctx, reaction)
}
