package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
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

	return &MongoOauthConnection{
		Id: oauthConnectionId,
		...
	}, nil
}

func (m MongoOauthConnectionMapper) MapToEntity(databaseUser *MongoOauthConnection) (*, error) {
	return &entity.OauthConnection{
		...
	}, nil
}
