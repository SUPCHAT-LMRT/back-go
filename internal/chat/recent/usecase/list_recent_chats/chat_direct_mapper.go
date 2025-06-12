package list_recent_chats

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
)

type DirectChatMapper struct {
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewDirectChatMapper(
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase,
) mapper.Mapper[*ChatDirectMapping, *ListRecentChatsUseCaseOutput] {
	return &DirectChatMapper{getUserByIdUseCase: getUserByIdUseCase}
}

func (g DirectChatMapper) MapFromEntity(recentChat *ListRecentChatsUseCaseOutput) (*ChatDirectMapping, error) {
	return nil, nil
}

func (g DirectChatMapper) MapToEntity(chatDirect *ChatDirectMapping) (*ListRecentChatsUseCaseOutput, error) {
	otherUserId := chatDirect.ChatDirect.User1Id
	if chatDirect.ChatDirect.User1Id == chatDirect.CurrentUserId {
		otherUserId = chatDirect.ChatDirect.User2Id
	}

	otherUser, err := g.getUserByIdUseCase.Execute(context.Background(), otherUserId)
	if err != nil {
		return nil, err
	}

	return &ListRecentChatsUseCaseOutput{
		Id:        entity.RecentChatId(otherUserId),
		Kind:      entity.RecentChatKindDirect,
		Name:      otherUser.FullName(),
		UpdatedAt: chatDirect.ChatDirect.UpdatedAt,
		LastMessage: &RecentChatLastMessage{
			Id:         entity.RecentChatId(chatDirect.LastMessageId),
			Content:    chatDirect.LastMessageContent,
			CreatedAt:  chatDirect.LastMessageCreatedAt,
			AuthorId:   chatDirect.LastMessageSenderId,
			AuthorName: otherUser.FullName(),
		},
	}, nil
}

type ChatDirectMapping struct {
	ChatDirect           *chat_direct_entity.ChatDirect
	CurrentUserId        user_entity.UserId
	LastMessageId        chat_direct_entity.ChatDirectId
	LastMessageContent   string
	LastMessageCreatedAt time.Time
	LastMessageSenderId  user_entity.UserId
}
