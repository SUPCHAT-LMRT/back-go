package list_recent_chats

import (
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type GroupMapper struct{}

func NewGroupMapper() mapper.Mapper[*GroupMapping, *ListRecentChatsUseCaseOutput] {
	return &GroupMapper{}
}

func (g GroupMapper) MapFromEntity(recentChat *ListRecentChatsUseCaseOutput) (*GroupMapping, error) {
	return nil, nil
}

func (g GroupMapper) MapToEntity(group *GroupMapping) (*ListRecentChatsUseCaseOutput, error) {
	return &ListRecentChatsUseCaseOutput{
		Id:        entity.RecentChatId(group.Group.Id),
		Kind:      entity.RecentChatKindGroup,
		Name:      group.Group.Name,
		UpdatedAt: group.Group.UpdatedAt,
		LastMessage: &RecentChatLastMessage{
			Id:         entity.RecentChatId(group.LastMessageId),
			Content:    group.LastMessageContent,
			CreatedAt:  group.LastMessageCreatedAt,
			AuthorId:   group.LastMessageSenderId,
			AuthorName: group.LastMessageSenderName,
		},
	}, nil
}

type GroupMapping struct {
	Group *group_entity.Group
	// TODO Set type to GroupMessageId
	LastMessageId         string
	LastMessageContent    string
	LastMessageCreatedAt  time.Time
	LastMessageSenderId   user_entity.UserId
	LastMessageSenderName string
}
