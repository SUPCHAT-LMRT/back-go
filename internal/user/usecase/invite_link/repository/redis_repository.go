package repository

import (
	"context"
	"errors"
	"time"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
)

var (
	buildInviteLinkRedisKey = func(token string) string {
		return "app_invite_link:" + token
	}
	buildInviteLinkEmailBindingRedisKey = func(email string) string {
		return "app_invite_link_email:" + email
	}
)

type RedisInviteLinkRepository struct {
	mapper mapper.Mapper[map[string]string, *entity.InviteLink]
	client *redis.Client
}

func NewRedisInviteLinkRepository(
	inviteLinkMapper mapper.Mapper[map[string]string, *entity.InviteLink],
	client *redis.Client,
) InviteLinkRepository {
	return &RedisInviteLinkRepository{mapper: inviteLinkMapper, client: client}
}

func (m RedisInviteLinkRepository) GenerateInviteLink(
	ctx context.Context,
	link *entity.InviteLink,
) error {
	link.ExpiresAt = time.Now().Add(RedisInviteLinkExpiredTime)

	databaseInviteLink, err := m.mapper.MapFromEntity(link)
	if err != nil {
		return err
	}

	_, err = m.client.Client.HSet(ctx, buildInviteLinkRedisKey(link.Token), databaseInviteLink).
		Result()
	if err != nil {
		return err
	}

	err = m.client.Client.Expire(ctx, buildInviteLinkRedisKey(link.Token), RedisInviteLinkExpiredTime).
		Err()
	if err != nil {
		return err
	}

	_, err = m.client.Client.Set(ctx, buildInviteLinkEmailBindingRedisKey(link.Email), link.Token, RedisInviteLinkExpiredTime).
		Result()
	if err != nil {
		return err
	}

	return nil
}

func (m RedisInviteLinkRepository) GetInviteLinkData(
	ctx context.Context,
	token string,
) (*entity.InviteLink, error) {
	databaseInviteLink, err := m.client.Client.HGetAll(ctx, buildInviteLinkRedisKey(token)).Result()
	if err != nil {
		return nil, err
	}

	if len(databaseInviteLink) == 0 {
		return nil, ErrInviteLinkNotFound
	}

	inviteLinkData, err := m.mapper.MapToEntity(databaseInviteLink)
	if err != nil {
		return nil, err
	}

	return inviteLinkData, nil
}

func (m RedisInviteLinkRepository) GetInviteLinkDataByEmail(
	ctx context.Context,
	email string,
) (*entity.InviteLink, error) {
	token, err := m.client.Client.Get(ctx, buildInviteLinkEmailBindingRedisKey(email)).Result()
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil, ErrInviteLinkNotFound
		}
		return nil, err
	}

	databaseInviteLink, err := m.client.Client.HGetAll(ctx, buildInviteLinkRedisKey(token)).Result()
	if err != nil {
		return nil, err
	}

	if len(databaseInviteLink) == 0 {
		return nil, ErrInviteLinkNotFound
	}

	inviteLinkData, err := m.mapper.MapToEntity(databaseInviteLink)
	if err != nil {
		return nil, err
	}

	return inviteLinkData, nil
}

func (m RedisInviteLinkRepository) DeleteInviteLink(ctx context.Context, token string) error {
	inviteLinkData, err := m.GetInviteLinkData(ctx, token)
	if err != nil {
		if errors.Is(err, ErrInviteLinkNotFound) {
			return nil // No need to delete if it doesn't exist
		}
	}

	_, err = m.client.Client.Del(ctx, buildInviteLinkRedisKey(token)).Result()
	if err != nil {
		return err
	}

	_, err = m.client.Client.Del(ctx, buildInviteLinkEmailBindingRedisKey(inviteLinkData.Email)).
		Result()
	if err != nil {
		return err
	}

	return nil
}

//nolint:revive
func (m RedisInviteLinkRepository) GetAllInviteLinks(
	ctx context.Context,
) ([]*entity.InviteLink, error) {
	keys, err := m.client.Client.Keys(ctx, "invite_link_workspace:*").Result()
	if err != nil {
		return nil, err
	}

	var inviteLinks []*entity.InviteLink
	for _, key := range keys {
		data, err := m.client.Client.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		inviteLink, err := m.mapper.MapToEntity(data)
		if err != nil {
			return nil, err
		}

		inviteLinks = append(inviteLinks, inviteLink)
	}

	return inviteLinks, nil
}
