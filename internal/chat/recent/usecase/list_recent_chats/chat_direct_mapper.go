package list_recent_chats

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
)

type DirectChatMapper struct {
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewDirectChatMapper(getUserByIdUseCase *get_by_id.GetUserByIdUseCase) mapper.Mapper[*ChatDirectMapping, *entity.RecentChat] {
	return &DirectChatMapper{getUserByIdUseCase: getUserByIdUseCase}
}

func (g DirectChatMapper) MapFromEntity(recentChat *entity.RecentChat) (*ChatDirectMapping, error) {
	return nil, nil
}

func (g DirectChatMapper) MapToEntity(chatDirect *ChatDirectMapping) (*entity.RecentChat, error) {
	otherUserId := chatDirect.ChatDirect.User1Id
	if chatDirect.ChatDirect.User1Id == chatDirect.CurrentUserId {
		otherUserId = chatDirect.ChatDirect.User2Id
	}

	otherUser, err := g.getUserByIdUseCase.Execute(context.Background(), otherUserId)
	if err != nil {
		return nil, err
	}

	return &entity.RecentChat{
		Id:        entity.RecentChatId(otherUserId),
		Kind:      entity.RecentChatKindDirect,
		Name:      otherUser.FullName(),
		UpdatedAt: chatDirect.ChatDirect.UpdatedAt,
	}, nil
}

type ChatDirectMapping struct {
	ChatDirect    *chat_direct_entity.ChatDirect
	CurrentUserId user_entity.UserId
}
