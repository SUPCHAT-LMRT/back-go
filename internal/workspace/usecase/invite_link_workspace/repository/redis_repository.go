package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/invite_link_workspace/entity"
	"time"
)

var (
	buildInviteLinkRedisKey = func(token string) string {
		return "invite_link_workspace:" + token
	}
)

type RedisInviteLinkRepository struct {
	mapper mapper.Mapper[map[string]string, *entity.InviteLink]
	client *redis.Client
}

func NewRedisInviteLinkRepository(mapper mapper.Mapper[map[string]string, *entity.InviteLink], client *redis.Client) InviteLinkRepository {
	return &RedisInviteLinkRepository{mapper: mapper, client: client}
}

func (m RedisInviteLinkRepository) GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error {

	link.ExpiresAt = time.Now().Add(RedisInviteLinkExpiredTime)

	databaseInviteLink, err := m.mapper.MapFromEntity(link)
	if err != nil {
		return err
	}

	_, err = m.client.Client.HSet(ctx, buildInviteLinkRedisKey(link.Token), databaseInviteLink).Result()
	if err != nil {
		return err
	}

	err = m.client.Client.Expire(ctx, buildInviteLinkRedisKey(link.Token), RedisInviteLinkExpiredTime).Err()
	if err != nil {
		return err
	}

	return nil

}

func (m RedisInviteLinkRepository) GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error) {

	databaseInviteLink, err := m.client.Client.HGetAll(ctx, buildInviteLinkRedisKey(token)).Result()
	if err != nil {
		return nil, err
	}

	inviteLinkData, err := m.mapper.MapToEntity(databaseInviteLink)
	if err != nil {
		return nil, err
	}

	return inviteLinkData, nil
}

func (m RedisInviteLinkRepository) DeleteInviteLink(ctx context.Context, token string) error {

	_, err := m.client.Client.Del(ctx, buildInviteLinkRedisKey(token)).Result()
	if err != nil {
		return err
	}

	return nil
}
