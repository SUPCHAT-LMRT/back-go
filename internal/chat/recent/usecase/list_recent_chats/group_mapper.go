package list_recent_chats

import (
	"context"
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"time"

	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type GroupMapper struct {
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewGroupMapper(getUserByIdUseCase *get_by_id.GetUserByIdUseCase) mapper.Mapper[*GroupMapping, *ListRecentChatsUseCaseOutput] {
	return &GroupMapper{getUserByIdUseCase: getUserByIdUseCase}
}

func (g GroupMapper) MapFromEntity(
	recentChat *ListRecentChatsUseCaseOutput,
) (*GroupMapping, error) {
	return nil, nil
}

func (g GroupMapper) MapToEntity(group *GroupMapping) (*ListRecentChatsUseCaseOutput, error) {
	lastMessage := RecentChatLastMessage{}

	if group.LastMessageId != "" {
		senderUser, err := g.getUserByIdUseCase.Execute(context.Background(), group.LastMessageSenderId)
		if err != nil {
			return nil, err
		}

		lastMessage = RecentChatLastMessage{
			Id:         entity.RecentChatId(group.LastMessageId),
			Content:    group.LastMessageContent,
			CreatedAt:  group.LastMessageCreatedAt,
			AuthorId:   group.LastMessageSenderId,
			AuthorName: senderUser.FullName(),
		}
	}

	return &ListRecentChatsUseCaseOutput{
		Id:          entity.RecentChatId(group.Group.Id),
		Kind:        entity.RecentChatKindGroup,
		Name:        group.Group.Name,
		UpdatedAt:   group.Group.UpdatedAt,
		LastMessage: &lastMessage,
	}, nil
}

type GroupMapping struct {
	Group                *group_entity.Group
	LastMessageId        group_chat_entity.GroupChatMessageId
	LastMessageContent   string
	LastMessageCreatedAt time.Time
	LastMessageSenderId  user_entity.UserId
}
