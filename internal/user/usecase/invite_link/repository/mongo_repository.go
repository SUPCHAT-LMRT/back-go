package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"time"
)

type MongoInviteLink struct {
	Token     string    `bson:"token"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type MongoInviteLinkRepository struct {
	mapper mapper.Mapper[*MongoInviteLink, *entity.InviteLink]
	client *mongo.Client
}

func NewMongoInviteLinkRepository(mapper mapper.Mapper[*MongoInviteLink, *entity.InviteLink], client *mongo.Client) InviteLinkRepository {
	return &MongoInviteLinkRepository{mapper: mapper, client: client}
}

func (m MongoInviteLinkRepository) GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error {
	databaseInviteLink, err := m.mapper.MapFromEntity(link)
	if err != nil {
		return err
	}

	_, err = m.client.Client.Database("supchat").Collection("invite_links").InsertOne(ctx, databaseInviteLink)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoInviteLinkRepository) GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error) {
	var inviteLink MongoInviteLink
	err := m.client.Client.Database("supchat").Collection("invite_links").FindOne(ctx, map[string]string{"token": token}).Decode(&inviteLink)
	if err != nil {
		return nil, err
	}

	inviteLinkData, err := m.mapper.MapToEntity(&inviteLink)
	if err != nil {
		return nil, err
	}

	return inviteLinkData, nil
}
