package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "workspaces_channels"
)

type MongoChannelRepositoryDeps struct {
	uberdig.In
	Client *mongo.Client
	Mapper mapper.Mapper[*MongoChannel, *entity.Channel]
}

type MongoChannelRepository struct {
	deps MongoChannelRepositoryDeps
}

type MongoChannel struct {
	Id          bson.ObjectID `bson:"_id"`
	Name        string        `bson:"name"`
	Topic       string        `bson:"topic"`
	WorkspaceId bson.ObjectID `bson:"workspace_id"`
}

func NewMongoChannelRepository(deps MongoChannelRepositoryDeps) ChannelRepository {
	return &MongoChannelRepository{deps: deps}
}

func (m MongoChannelRepository) Create(ctx context.Context, channel *entity.Channel) error {
	channel.Id = entity.ChannelId(bson.NewObjectID().Hex())

	mongoChannel, err := m.deps.Mapper.MapFromEntity(channel)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoChannel)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChannelRepository) GetById(ctx context.Context, id entity.ChannelId) (*entity.Channel, error) {
	objectId, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var mongoChannel MongoChannel
	err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": objectId}).Decode(&mongoChannel)
	if err != nil {
		return nil, err
	}

	channel, err := m.deps.Mapper.MapToEntity(&mongoChannel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (m MongoChannelRepository) List(ctx context.Context, workspaceId workspace_entity.WorkspaceId) ([]*entity.Channel, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(workspaceId))
	if err != nil {
		return nil, err
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{"workspace_id": workspaceObjectId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoChannels []*MongoChannel
	err = cursor.All(ctx, &mongoChannels)
	if err != nil {
		return nil, err
	}

	channels := make([]*entity.Channel, len(mongoChannels))
	for i, mongoChannel := range mongoChannels {
		channel, err := m.deps.Mapper.MapToEntity(mongoChannel)
		if err != nil {
			return nil, err
		}
		channels[i] = channel
	}

	return channels, nil
}
