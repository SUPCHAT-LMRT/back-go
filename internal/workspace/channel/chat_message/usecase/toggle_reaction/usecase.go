package toggle_reaction

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
)

type ToggleReactionUseCase struct {
	repository repository.ChannelMessageRepository
}

func NewToggleReactionUseCase(repository repository.ChannelMessageRepository) *ToggleReactionUseCase {
	return &ToggleReactionUseCase{repository: repository}
}

func (u *ToggleReactionUseCase) Execute(ctx context.Context, messageId entity.ChannelMessageId, userId user_entity.UserId, reaction string) (added bool, err error) {
	return u.repository.ToggleReaction(ctx, messageId, userId, reaction)
}
