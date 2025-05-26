package list_messages

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_chat_message_repository "github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type ListGroupChatMessagesUseCase struct {
	repository group_chat_message_repository.GroupChatMessageRepository
}

func NewListGroupChatMessagesUseCase(
	repository group_chat_message_repository.GroupChatMessageRepository,
) *ListGroupChatMessagesUseCase {
	return &ListGroupChatMessagesUseCase{repository: repository}
}

func (u ListGroupChatMessagesUseCase) Execute(
	ctx context.Context,
	groupId group_entity.GroupId,
) ([]*entity.GroupChatMessage, error) {
	return u.repository.ListByGroupId(ctx, groupId)
}
