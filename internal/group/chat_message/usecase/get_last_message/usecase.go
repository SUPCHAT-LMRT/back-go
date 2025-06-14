package get_last_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"

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
) (*LastMessageResponse, error) {
	lastMessage, err := u.repository.GetLastMessage(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if lastMessage == nil {
		return nil, nil // No messages found
	}

	return &LastMessageResponse{
		Id:        lastMessage.Id,
		GroupId:   lastMessage.GroupId,
		Content:   lastMessage.Content,
		CreatedAt: lastMessage.CreatedAt,
		AuthorId:  lastMessage.AuthorId,
	}, err
}

type LastMessageResponse struct {
	Id        entity.GroupChatMessageId
	GroupId   group_entity.GroupId
	Content   string
	CreatedAt time.Time
	AuthorId  user_entity.UserId
}
