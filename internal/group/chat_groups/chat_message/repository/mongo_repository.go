package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_groups/chat_message/entity"
	group_chat_message_entity "github.com/supchat-lmrt/back-go/internal/group/chat_groups/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"time"
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

type MongoChatMessageRepository struct {
	deps MongoChatMessageRepositoryDeps
}

type MongoGroupChatMessage struct {
	Id        bson.ObjectID `bson:"_id"`
	GroupId   bson.ObjectID `bson:"group_id"`
	AuthorId  bson.ObjectID `bson:"user_id"`
	Content   string        `bson:"content"`
	CreatedAt time.Time     `bson:"created_at"`
}

func NewMongoChatMessageRepository(deps MongoChatMessageRepositoryDeps) GroupChatMessageRepository {
	return &MongoChatMessageRepository{deps: deps}
}

func (m MongoChatMessageRepository) Create(ctx context.Context, message *entity.GroupChatMessage) error {
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

func (m MongoChatMessageRepository) ListByGroupId(ctx context.Context, groupId group_chat_message_entity.ChatGroupId) ([]*entity.GroupChatMessage, error) {
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
