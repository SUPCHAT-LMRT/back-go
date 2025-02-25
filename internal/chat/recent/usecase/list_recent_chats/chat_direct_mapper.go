package list_recent_chats

import (
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
)

type DirectChatMapper struct {
	getUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

func NewDirectChatMapper(getUserByIdUseCase *get_by_id.GetUserByIdUseCase) mapper.Mapper[*chat_direct_entity.ChatDirect, *entity.RecentChat] {
	return &DirectChatMapper{getUserByIdUseCase: getUserByIdUseCase}
}

func (g DirectChatMapper) MapFromEntity(recentChat *entity.RecentChat) (*chat_direct_entity.ChatDirect, error) {
	return nil, nil
}

func (g DirectChatMapper) MapToEntity(chatDirect *chat_direct_entity.ChatDirect) (*entity.RecentChat, error) {

	return &entity.RecentChat{
		Id:        entity.RecentChatId(chatDirect.Id),
		Kind:      entity.RecentChatKindDirect,
		AvatarUrl: "http://localhost:9000/users-avatars/" + chatDirect.User2Id.String(),
		UpdatedAt: chatDirect.UpdatedAt,
	}, nil
}
