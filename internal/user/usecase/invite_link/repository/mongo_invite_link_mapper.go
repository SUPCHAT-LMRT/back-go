package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
)

type MongoInviteLinkMapper struct{}

func NewMongoInviteLinkMapper() mapper.Mapper[*MongoInviteLink, *entity.InviteLink] {
	return &MongoInviteLinkMapper{}
}

func (m MongoInviteLinkMapper) MapFromEntity(entity *entity.InviteLink) (*MongoInviteLink, error) {
	return &MongoInviteLink{
		Token:     entity.Token,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
		ExpiresAt: entity.ExpiresAt,
	}, nil
}

func (m MongoInviteLinkMapper) MapToEntity(databaseInviteLink *MongoInviteLink) (*entity.InviteLink, error) {
	return &entity.InviteLink{
		Token:     databaseInviteLink.Token,
		FirstName: databaseInviteLink.FirstName,
		LastName:  databaseInviteLink.LastName,
		Email:     databaseInviteLink.Email,
		ExpiresAt: databaseInviteLink.ExpiresAt,
	}, nil
}
