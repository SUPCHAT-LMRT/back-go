package toggle_reaction

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ToggleReactionDirectMessageUseCase struct {
	repository repository.ChatDirectRepository
}

func NewToggleReactionDirectMessageUseCase(repository repository.ChatDirectRepository) *ToggleReactionDirectMessageUseCase {
	return &ToggleReactionDirectMessageUseCase{repository: repository}
}

func (u *ToggleReactionDirectMessageUseCase) Execute(ctx context.Context, messageId entity.ChatDirectId, userId user_entity.UserId, reaction string) (added bool, err error) {
	return u.repository.ToggleReaction(ctx, messageId, userId, reaction)
}
