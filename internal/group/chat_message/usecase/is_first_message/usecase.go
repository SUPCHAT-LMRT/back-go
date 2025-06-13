package is_first_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type IsFirstGroupChatMessageUseCase struct {
	repository repository.ChatMessageRepository
}

func NewIsFirstGroupChatMessageUseCase(
	repository repository.ChatMessageRepository,
) *IsFirstGroupChatMessageUseCase {
	return &IsFirstGroupChatMessageUseCase{repository: repository}
}

func (u *IsFirstGroupChatMessageUseCase) Execute(
	ctx context.Context,
	groupId group_entity.GroupId,
) (bool, error) {
	return u.repository.IsFirstMessage(ctx, groupId)
}
