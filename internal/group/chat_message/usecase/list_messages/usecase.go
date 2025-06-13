package list_messages

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type QueryParams struct {
	Limit           int
	Before          time.Time
	After           time.Time
	AroundMessageId entity.GroupChatMessageId
}

type ListGroupChatMessagesUseCase struct {
	repository repository.ChatMessageRepository
}

func NewListGroupChatMessagesUseCase(
	repository repository.ChatMessageRepository,
) *ListGroupChatMessagesUseCase {
	return &ListGroupChatMessagesUseCase{repository: repository}
}

func (u *ListGroupChatMessagesUseCase) Execute(
	ctx context.Context,
	groupId group_entity.GroupId,
	params QueryParams,
) ([]*entity.GroupChatMessage, error) {
	return u.repository.ListMessages(ctx, groupId, repository.ListMessagesQueryParams{
		Limit:           params.Limit,
		Before:          params.Before,
		After:           params.After,
		AroundMessageId: params.AroundMessageId,
	})
}
