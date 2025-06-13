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
	return &ListRecentChatsUseCaseOutput{
		Id:   response.Id,
		Kind: response.Kind,
		Name: response.Name,
		LastMessage: &RecentChatLastMessage{
			Id:         response.LastMessage.Id,
			Content:    response.LastMessage.Content,
			CreatedAt:  response.LastMessage.CreatedAt,
			AuthorId:   response.LastMessage.AuthorId,
			AuthorName: response.LastMessage.AuthorName,
		},
	}, nil
}

func (r ResponseMapper) MapToEntity(
	entityRecentChat *ListRecentChatsUseCaseOutput,
) (*RecentChatResponse, error) {
	return &RecentChatResponse{
		Id:   entityRecentChat.Id,
		Kind: entityRecentChat.Kind,
		Name: entityRecentChat.Name,
		LastMessage: RecentChatLastMessageResponse{
			Id:         entityRecentChat.LastMessage.Id,
			Content:    entityRecentChat.LastMessage.Content,
			CreatedAt:  entityRecentChat.LastMessage.CreatedAt,
			AuthorId:   entityRecentChat.LastMessage.AuthorId,
			AuthorName: entityRecentChat.LastMessage.AuthorName,
		},
	}, nil
}
