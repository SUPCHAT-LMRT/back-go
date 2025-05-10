package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoChannelMapper struct{}

func NewMongoChannelMapper() mapper.Mapper[*MongoChannel, *entity.Channel] {
	return &MongoChannelMapper{}
}

func (m MongoChannelMapper) MapFromEntity(entity *entity.Channel) (*MongoChannel, error) {
	channelObjectId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		return nil, err
	}

	workspaceObjectId, err := bson.ObjectIDFromHex(string(entity.WorkspaceId))
	if err != nil {
		return nil, err
	}

	return &MongoChannel{
		Id:          channelObjectId,
		Name:        entity.Name,
		Topic:       entity.Topic,
		Kind:        entity.Kind.String(),
		IsPrivate:   entity.IsPrivate,
		Members:     entity.Members,
		WorkspaceId: workspaceObjectId,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Index:       entity.Index,
	}, nil
}

func (m MongoChannelMapper) MapToEntity(mongo *MongoChannel) (*entity.Channel, error) {
	return &entity.Channel{
		Id:          entity.ChannelId(mongo.Id.Hex()),
		Name:        mongo.Name,
		Topic:       mongo.Topic,
		Kind:        entity.ChannelKind(mongo.Kind),
		IsPrivate:   mongo.IsPrivate,
		Members:     mongo.Members,
		WorkspaceId: workspace_entity.WorkspaceId(mongo.WorkspaceId.Hex()),
		CreatedAt:   mongo.CreatedAt,
		UpdatedAt:   mongo.UpdatedAt,
		Index:       mongo.Index,
	}, nil
}
