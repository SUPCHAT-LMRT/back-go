package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/entity"
	"time"
)

type RedisInviteLinkWorkspaceMapper struct{}

func NewRedisInviteLinkMapper() mapper.Mapper[map[string]string, *entity.InviteLink] {
	return &RedisInviteLinkWorkspaceMapper{}
}

func (m RedisInviteLinkWorkspaceMapper) MapFromEntity(entity *entity.InviteLink) (map[string]string, error) {
	return map[string]string{
		"token":       entity.Token,
		"workspaceId": entity.WorkspaceId.String(),
		"expires_at":  entity.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (m RedisInviteLinkWorkspaceMapper) MapToEntity(databaseInviteLink map[string]string) (*entity.InviteLink, error) {
	expiresAt, err := time.Parse(time.RFC3339, databaseInviteLink["expires_at"])
	if err != nil {
		return nil, err
	}

	return &entity.InviteLink{
		Token:       databaseInviteLink["token"],
		WorkspaceId: workspace_entity.WorkspaceId(databaseInviteLink["workspaceId"]),
		ExpiresAt:   expiresAt,
	}, nil
}
