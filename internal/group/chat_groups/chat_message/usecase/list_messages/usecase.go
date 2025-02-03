package list_messages

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_groups/chat_message/entity"
	group_chat_message_repository "github.com/supchat-lmrt/back-go/internal/group/chat_groups/chat_message/repository"
	group_chat_message_entity "github.com/supchat-lmrt/back-go/internal/group/chat_groups/entity"
)

type ListGroupChatMessagesUseCase struct {
	repository group_chat_message_repository.GroupChatMessageRepository
}

func NewListGroupChatMessagesUseCase(repository group_chat_message_repository.GroupChatMessageRepository) *ListGroupChatMessagesUseCase {
	return &ListGroupChatMessagesUseCase{repository: repository}
}

func (u ListGroupChatMessagesUseCase) Execute(ctx context.Context, groupId group_chat_message_entity.ChatGroupId) ([]*entity.GroupChatMessage, error) {
	return u.repository.ListByGroupId(ctx, groupId)
}
