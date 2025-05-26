package repository

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
)

type RedisInviteLinkMapper struct{}

func NewRedisInviteLinkMapper() mapper.Mapper[map[string]string, *entity.InviteLink] {
	return &RedisInviteLinkMapper{}
}

func (m RedisInviteLinkMapper) MapFromEntity(
	entityInviteLink *entity.InviteLink,
) (map[string]string, error) {
	return map[string]string{
		"token":      entityInviteLink.Token,
		"first_name": entityInviteLink.FirstName,
		"last_name":  entityInviteLink.LastName,
		"email":      entityInviteLink.Email,
		"expires_at": entityInviteLink.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (m RedisInviteLinkMapper) MapToEntity(
	databaseInviteLink map[string]string,
) (*entity.InviteLink, error) {
	expiresAt, err := time.Parse(time.RFC3339, databaseInviteLink["expires_at"])
	if err != nil {
		return nil, err
	}

	return &entity.InviteLink{
		Token:     databaseInviteLink["token"],
		FirstName: databaseInviteLink["first_name"],
		LastName:  databaseInviteLink["last_name"],
		Email:     databaseInviteLink["email"],
		ExpiresAt: expiresAt,
	}, nil
}
