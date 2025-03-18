package list_recent_chats

import (
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
)

type ResponseMapper struct{}

func NewResponseMapper() mapper.Mapper[*entity.RecentChat, *RecentChatResponse] {
	return &ResponseMapper{}
}

func (r ResponseMapper) MapFromEntity(response *RecentChatResponse) (*entity.RecentChat, error) {
	return &entity.RecentChat{
		Id:   response.Id,
		Kind: response.Kind,
		Name: response.Name,
	}, nil
}

func (r ResponseMapper) MapToEntity(entity *entity.RecentChat) (*RecentChatResponse, error) {
	return &RecentChatResponse{
		Id:   entity.Id,
		Kind: entity.Kind,
		Name: entity.Name,
	}, nil
}
