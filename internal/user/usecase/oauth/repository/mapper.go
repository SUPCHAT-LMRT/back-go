package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoOauthConnectionMapper struct{}

func NewMongoOauthConnectionMapper() mapper.Mapper[*MongoOauthConnection, *entity.OauthConnection] {
	return &MongoOauthConnectionMapper{}
}

func (m MongoOauthConnectionMapper) MapFromEntity(entity *entity.OauthConnection) (*MongoOauthConnection, error) {
	oauthConnectionId, err := bson.ObjectIDFromHex(entity.Id.String())
	if err != nil {
		return nil, err
	}

	userId, err := bson.ObjectIDFromHex(entity.UserId.String())
	if err != nil {
		return nil, err
	}

	return &MongoOauthConnection{
		Id:          oauthConnectionId,
		UserId:      userId,
		Provider:    entity.Provider,
		OauthEmail:  entity.OauthEmail,
		OauthUserId: entity.OauthUserId,
	}, nil
}

func (m MongoOauthConnectionMapper) MapToEntity(databaseUser *MongoOauthConnection) (*entity.OauthConnection, error) {
	return &entity.OauthConnection{
		Id:          entity.OauthConnectionId(databaseUser.Id.Hex()),
		UserId:      user_entity.UserId(databaseUser.UserId.Hex()),
		Provider:    databaseUser.Provider,
		OauthEmail:  databaseUser.OauthEmail,
		OauthUserId: databaseUser.OauthUserId,
	}, nil
}
