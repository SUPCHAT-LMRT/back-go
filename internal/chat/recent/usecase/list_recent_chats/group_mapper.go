package list_recent_chats

import (
	"github.com/supchat-lmrt/back-go/internal/chat/recent/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
)

type GroupMapper struct{}

func NewGroupMapper() mapper.Mapper[*group_entity.Group, *entity.RecentChat] {
	return &GroupMapper{}
}

func (g GroupMapper) MapFromEntity(recentChat *entity.RecentChat) (*group_entity.Group, error) {
	return nil, nil
}

func (g GroupMapper) MapToEntity(group *group_entity.Group) (*entity.RecentChat, error) {
	return &entity.RecentChat{
		Id:        entity.RecentChatId(group.Id),
		Kind:      entity.RecentChatKindGroup,
		Name:      group.Name,
		UpdatedAt: group.UpdatedAt,
	}, nil
}
