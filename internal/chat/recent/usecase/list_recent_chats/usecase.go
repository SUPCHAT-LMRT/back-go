package list_recent_chats

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/get_last_message"
	get_last_message2 "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/get_last_message"
	"sort"
	"time"

	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_recent_groups"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type ListRecentChatsUseCaseDeps struct {
	uberdig.In
	ListRecentGroupsUseCase         *list_recent_groups.ListRecentGroupsUseCase
	GetLastDirectChatMessageUseCase *get_last_message2.GetLastDirectChatMessageUseCase
	ListRecentChatDirectUseCase     *list_recent_direct_chats.ListRecentChatDirectUseCase
	GetLastGroupChatMessageUseCase  *get_last_message.GetLastGroupChatMessageUseCase
	GroupMapper                     mapper.Mapper[*GroupMapping, *ListRecentChatsUseCaseOutput]
	DirectMapper                    mapper.Mapper[*ChatDirectMapping, *ListRecentChatsUseCaseOutput]
}

type ListRecentChatsUseCase struct {
	deps ListRecentChatsUseCaseDeps
}

func NewListRecentChatsUseCase(deps ListRecentChatsUseCaseDeps) *ListRecentChatsUseCase {
	return &ListRecentChatsUseCase{deps: deps}
}

//nolint:revive
func (u *ListRecentChatsUseCase) Execute(
	ctx context.Context,
	currentUserId user_entity.UserId,
) ([]*ListRecentChatsUseCaseOutput, error) {
	groups, err := u.deps.ListRecentGroupsUseCase.Execute(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent groups: %w", err)
	}

	directs, err := u.deps.ListRecentChatDirectUseCase.Execute(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent chat directs: %w", err)
	}

	var recentChats []*ListRecentChatsUseCaseOutput

	for _, group := range groups {
		lastMessage, err := u.deps.GetLastGroupChatMessageUseCase.Execute(ctx, group.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get last group chat message: %w", err)
		}

		mapping := GroupMapping{
			Group: group,
		}
		if lastMessage != nil {
			mapping.LastMessageId = lastMessage.Id
			mapping.LastMessageContent = lastMessage.Content
			mapping.LastMessageCreatedAt = lastMessage.CreatedAt
			mapping.LastMessageSenderId = lastMessage.AuthorId
		}

		fromEntity, err := u.deps.GroupMapper.MapToEntity(&mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to map group to entity: %w", err)
		}

		recentChats = append(recentChats, fromEntity)
	}

	for _, direct := range directs {
		lastMessage, err := u.deps.GetLastDirectChatMessageUseCase.Execute(
			ctx,
			direct.User1Id,
			direct.User2Id,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get last direct chat message: %w", err)
		}

		mapping := ChatDirectMapping{
			ChatDirect:    direct,
			CurrentUserId: currentUserId,
		}
		if lastMessage != nil {
			mapping.LastMessageId = lastMessage.ChatId
			mapping.LastMessageContent = lastMessage.Content
			mapping.LastMessageCreatedAt = lastMessage.CreatedAt
			mapping.LastMessageSenderId = lastMessage.AuthorId
		}

		fromEntity, err := u.deps.DirectMapper.MapToEntity(
			&mapping,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to map direct to entity: %w", err)
		}

		recentChats = append(recentChats, fromEntity)
	}

	// Sort by updated at
	sort.Slice(recentChats, func(i, j int) bool {
		return recentChats[i].UpdatedAt.After(recentChats[j].UpdatedAt)
	})

	return recentChats, nil
}

type ListRecentChatsUseCaseOutput struct {
	Id          entity.RecentChatId
	Kind        entity.RecentChatKind
	Name        string
	UpdatedAt   time.Time
	LastMessage *RecentChatLastMessage
}

type RecentChatLastMessage struct {
	Id         entity.RecentChatId
	Content    string
	CreatedAt  time.Time
	AuthorId   user_entity.UserId
	AuthorName string
}
