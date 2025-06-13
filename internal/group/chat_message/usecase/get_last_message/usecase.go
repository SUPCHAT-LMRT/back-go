package get_last_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type GetLastGroupChatMessageUseCase struct {
	repository repository.ChatMessageRepository
}

func NewGetLastGroupChatMessageUseCase(
	repository repository.ChatMessageRepository,
) *GetLastGroupChatMessageUseCase {
	return &GetLastGroupChatMessageUseCase{repository: repository}
}

func (u *GetLastGroupChatMessageUseCase) Execute(
	ctx context.Context,
	groupId group_entity.GroupId,
) (*entity.GroupChatMessage, error) {
	return u.repository.GetLastMessage(ctx, groupId)
}
