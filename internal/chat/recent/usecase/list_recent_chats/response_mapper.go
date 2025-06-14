package list_recent_chats

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
)

type ResponseMapper struct{}

func NewResponseMapper() mapper.Mapper[*ListRecentChatsUseCaseOutput, *RecentChatResponse] {
	return &ResponseMapper{}
}

func (r ResponseMapper) MapFromEntity(
	response *RecentChatResponse,
) (*ListRecentChatsUseCaseOutput, error) {
	lastMessage := RecentChatLastMessage{}

	if response.LastMessage != nil {
		lastMessage = RecentChatLastMessage{
			Id:         response.LastMessage.Id,
			Content:    response.LastMessage.Content,
			CreatedAt:  response.LastMessage.CreatedAt,
			AuthorId:   response.LastMessage.AuthorId,
			AuthorName: response.LastMessage.AuthorName,
		}
	}

	return &ListRecentChatsUseCaseOutput{
		Id:          response.Id,
		Kind:        response.Kind,
		Name:        response.Name,
		LastMessage: &lastMessage,
	}, nil
}

func (r ResponseMapper) MapToEntity(
	entityRecentChat *ListRecentChatsUseCaseOutput,
) (*RecentChatResponse, error) {
	lastMessage := RecentChatLastMessageResponse{}

	if entityRecentChat.LastMessage != nil {
		lastMessage = RecentChatLastMessageResponse{
			Id:         entityRecentChat.LastMessage.Id,
			Content:    entityRecentChat.LastMessage.Content,
			CreatedAt:  entityRecentChat.LastMessage.CreatedAt,
			AuthorId:   entityRecentChat.LastMessage.AuthorId,
			AuthorName: entityRecentChat.LastMessage.AuthorName,
		}
	}

	return &RecentChatResponse{
		Id:          entityRecentChat.Id,
		Kind:        entityRecentChat.Kind,
		Name:        entityRecentChat.Name,
		LastMessage: &lastMessage,
	}, nil
}
