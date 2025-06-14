package toggle_reaction

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ToggleReactionInput struct {
	MessageId entity.GroupChatMessageId
	UserId    user_entity.UserId
	Reaction  string
}

type ToggleGroupChatReactionUseCase struct {
	repository repository.ChatMessageRepository
}

func NewToggleGroupChatReactionUseCase(
	repository repository.ChatMessageRepository,
) *ToggleGroupChatReactionUseCase {
	return &ToggleGroupChatReactionUseCase{repository: repository}
}

func (u *ToggleGroupChatReactionUseCase) Execute(
	ctx context.Context,
	input ToggleReactionInput,
) (bool, error) {
	return u.repository.ToggleReaction(ctx, input.MessageId, input.UserId, input.Reaction)
}
