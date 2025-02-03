package list_recent_chats

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_recent_groups"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	uberdig "go.uber.org/dig"
	"sort"
)

type ListRecentChatsUseCaseDeps struct {
	uberdig.In
	ListRecentGroupsUseCase     *list_recent_groups.ListRecentGroupsUseCase
	ListRecentChatDirectUseCase *list_recent_direct_chats.ListRecentChatDirectUseCase
	GroupMapper                 mapper.Mapper[*group_entity.Group, *entity.RecentChat]
	DirectMapper                mapper.Mapper[*chat_direct_entity.ChatDirect, *entity.RecentChat]
}

type ListRecentChatsUseCase struct {
	deps ListRecentChatsUseCaseDeps
}

func NewListRecentChatsUseCase(deps ListRecentChatsUseCaseDeps) *ListRecentChatsUseCase {
	return &ListRecentChatsUseCase{deps: deps}
}

func (u *ListRecentChatsUseCase) Execute(ctx context.Context) ([]*entity.RecentChat, error) {
	// Call the ListRecentGroupsUseCase and ListRecentChatDirectUseCase and sort by updated_at

	groups, err := u.deps.ListRecentGroupsUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}

	directs, err := u.deps.ListRecentChatDirectUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}

	// Sort by updated_at
	var recentChats []*entity.RecentChat

	for _, group := range groups {
		fromEntity, err := u.deps.GroupMapper.MapToEntity(group)
		if err != nil {
			return nil, err
		}

		recentChats = append(recentChats, fromEntity)
	}

	for _, direct := range directs {
		fromEntity, err := u.deps.DirectMapper.MapToEntity(direct)
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
