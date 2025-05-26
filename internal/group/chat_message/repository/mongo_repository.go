package repository

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "group_messages"
)

type MongoChatMessageRepositoryDeps struct {
	uberdig.In
	Client *mongo.Client
	Mapper mapper.Mapper[*MongoGroupChatMessage, *entity.GroupChatMessage]
}

type MongoGroupChatMessageRepository struct {
	deps MongoChatMessageRepositoryDeps
}

type MongoGroupChatMessage struct {
	Id        bson.ObjectID `bson:"_id"`
	GroupId   bson.ObjectID `bson:"group_id"`
	AuthorId  bson.ObjectID `bson:"user_id"`
	Content   string        `bson:"content"`
	CreatedAt time.Time     `bson:"created_at"`
}

func NewMongoGroupChatMessageRepository(
	deps MongoChatMessageRepositoryDeps,
) GroupChatMessageRepository {
	return &MongoGroupChatMessageRepository{deps: deps}
}

func (m MongoGroupChatMessageRepository) Create(
	ctx context.Context,
	message *entity.GroupChatMessage,
) error {
	message.Id = entity.GroupChatMessageId(bson.NewObjectID().Hex())
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	mongoMessage, err := m.deps.Mapper.MapFromEntity(message)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, mongoMessage)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoGroupChatMessageRepository) ListByGroupId(
	ctx context.Context,
	groupId group_entity.GroupId,
) ([]*entity.GroupChatMessage, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	groupObjectId, err := bson.ObjectIDFromHex(string(groupId))
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"group_id": groupObjectId})
	if err != nil {
		return nil, err
	}

	var messages []*entity.GroupChatMessage
	for cursor.Next(ctx) {
		var mongoMessage MongoGroupChatMessage
		err = cursor.Decode(&mongoMessage)
		if err != nil {
			return nil, err
		}

		message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
