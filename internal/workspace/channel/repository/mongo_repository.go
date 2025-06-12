package repository

import (
	"context"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"time"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	Kind        string        `bson:"kind"`
	IsPrivate   bool          `bson:"is_private"`
	Members     []string      `bson:"members"`
	WorkspaceId bson.ObjectID `bson:"workspace_id"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	Index       int           `bson:"index"`
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

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoChannel)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChannelRepository) GetById(
	ctx context.Context,
	id entity.ChannelId,
) (*entity.Channel, error) {
	objectId, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var mongoChannel MongoChannel
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": objectId}).
		Decode(&mongoChannel)
	if err != nil {
		return nil, err
	}

	channel, err := m.deps.Mapper.MapToEntity(&mongoChannel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (m MongoChannelRepository) List(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
) ([]*entity.Channel, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(workspaceId))
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"workspace_id": workspaceObjectId,
		"is_private":   false,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "index", Value: 1}})

	cursor, err := m.deps.Client.Client.
		Database(databaseName).
		Collection(collectionName).
		Find(ctx, filter, findOptions)
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

func (m MongoChannelRepository) CountByWorkspaceId(
	ctx context.Context,
	id workspace_entity.WorkspaceId,
) (uint, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		return 0, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"workspace_id": workspaceObjectId})
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (m MongoChannelRepository) UpdateIndex(
	ctx context.Context,
	channelId entity.ChannelId,
	index int,
) error {
	objectId, err := bson.ObjectIDFromHex(string(channelId))
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"index":      index,
			"updated_at": time.Now(),
		},
	}

	_, err = m.deps.Client.Client.
		Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": objectId}, update)

	return err
}

func (m MongoChannelRepository) Delete(ctx context.Context, channelId entity.ChannelId) error {
	objectId, err := bson.ObjectIDFromHex(string(channelId))
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.
		Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{"_id": objectId})

	return err
}

func (m MongoChannelRepository) ListPrivateChannelsByUser(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
	memberId workspace_member_entity.WorkspaceMemberId,
) ([]*entity.Channel, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(workspaceId))
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"workspace_id": workspaceObjectId,
		"is_private":   true,
		"members":      bson.M{"$in": []string{string(memberId)}},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "index", Value: 1}})

	cursor, err := m.deps.Client.Client.
		Database(databaseName).
		Collection(collectionName).
		Find(ctx, filter, findOptions)
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

func (m MongoChannelRepository) ListMembersOfPrivateChannel(
	ctx context.Context,
	channelId entity.ChannelId,
) ([]workspace_member_entity.WorkspaceMemberId, error) {
	objectId, err := bson.ObjectIDFromHex(string(channelId))
	if err != nil {
		return nil, err
	}

	var mongoChannel MongoChannel
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": objectId, "is_private": true}).
		Decode(&mongoChannel)
	if err != nil {
		return nil, err
	}

	members := make([]workspace_member_entity.WorkspaceMemberId, len(mongoChannel.Members))
	for i, memberId := range mongoChannel.Members {
		members[i] = workspace_member_entity.WorkspaceMemberId(memberId)
	}
	return members, nil
}
