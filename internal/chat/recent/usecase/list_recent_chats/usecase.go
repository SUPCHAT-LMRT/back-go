package list_recent_chats

import (
	"context"
	"sort"
	"time"

	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_recent_groups"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/get_last_message"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type ListRecentChatsUseCaseDeps struct {
	uberdig.In
	ListRecentGroupsUseCase         *list_recent_groups.ListRecentGroupsUseCase
	ListRecentChatDirectUseCase     *list_recent_direct_chats.ListRecentChatDirectUseCase
	GetLastDirectChatMessageUseCase *get_last_message.GetLastDirectChatMessageUseCase
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
	// Call the ListRecentGroupsUseCase and ListRecentChatDirectUseCase and sort by updated_at

	groups, err := u.deps.ListRecentGroupsUseCase.Execute(ctx, currentUserId)
	if err != nil {
		return nil, err
	}

	directs, err := u.deps.ListRecentChatDirectUseCase.Execute(ctx, currentUserId)
	if err != nil {
		return nil, err
	}

	// Sort by updated_at
	var recentChats []*ListRecentChatsUseCaseOutput

	for _, group := range groups {
		// TODO implement last message for groups
		// lastMessage, err := u.deps.GetLastDirectChatMessageUseCase.Execute(ctx, direct.User1Id, direct.User2Id)
		// if err != nil {
		// 	return nil, err
		// }

		fromEntity, err := u.deps.GroupMapper.MapToEntity(&GroupMapping{
			Group: group,
			// LastMessageContent:   lastMessage.Content,
			// LastMessageCreatedAt: lastMessage.CreatedAt,
			// LastMessageSenderId:  lastMessage.AuthorId,
		})
		if err != nil {
			return nil, err
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
			return nil, err
		}

		fromEntity, err := u.deps.DirectMapper.MapToEntity(
			&ChatDirectMapping{
				ChatDirect:           direct,
				CurrentUserId:        currentUserId,
				LastMessageId:        lastMessage.ChatId,
				LastMessageContent:   lastMessage.Content,
				LastMessageCreatedAt: lastMessage.CreatedAt,
				LastMessageSenderId:  lastMessage.AuthorId,
			},
		)
		if err != nil {
			return nil, err
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
