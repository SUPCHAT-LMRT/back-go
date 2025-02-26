package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	databaseName   = "supchat"
	collectionName = "channel_messages"
)

type MongoChannelMessageRepositoryDeps struct {
	uberdig.In
	Client *mongo.Client
	Mapper mapper.Mapper[*MongoChannelMessage, *entity.ChannelMessage]
}

type MongoChannelMessageRepository struct {
	deps MongoChannelMessageRepositoryDeps
}

type MongoChannelMessage struct {
	Id        bson.ObjectID                  `bson:"_id"`
	ChannelId bson.ObjectID                  `bson:"channel_id"`
	AuthorId  bson.ObjectID                  `bson:"user_id"`
	Content   string                         `bson:"content"`
	CreatedAt time.Time                      `bson:"created_at"`
	Reactions []*MongoChannelMessageReaction `bson:"reactions"`
}

type MongoChannelMessageReaction struct {
	Id       bson.ObjectID `bson:"_id"`
	UserId   bson.ObjectID `bson:"user_id"`
	Reaction string        `bson:"reaction"`
}

func NewMongoChannelMessageRepository(deps MongoChannelMessageRepositoryDeps) ChannelMessageRepository {
	return &MongoChannelMessageRepository{deps: deps}
}

func (m MongoChannelMessageRepository) Create(ctx context.Context, message *entity.ChannelMessage) error {
	message.Id = entity.ChannelMessageId(bson.NewObjectID().Hex())
	message.CreatedAt = time.Now()

	mongoMessage, err := m.deps.Mapper.MapFromEntity(message)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoMessage)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChannelMessageRepository) ListByChannelId(ctx context.Context, channelId channel_entity.ChannelId) ([]*entity.ChannelMessage, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	channelObjectId, err := bson.ObjectIDFromHex(string(channelId))
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"channel_id": channelObjectId})
	if err != nil {
		return nil, err
	}

	var messages []*entity.ChannelMessage
	for cursor.Next(ctx) {
		var mongoMessage MongoChannelMessage
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

func (m MongoChannelMessageRepository) AddReaction(ctx context.Context, reaction entity.ChannelMessageReaction) error {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(reaction.MessageId.String())
	if err != nil {
		return err
	}

	userObjectId, err := bson.ObjectIDFromHex(reaction.UserId.String())
	if err != nil {
		return err
	}

	reactionObjectId, err := bson.ObjectIDFromHex(reaction.Id.String())
	if err != nil {
		return err
	}

	mongoReaction := MongoChannelMessageReaction{
		Id:       reactionObjectId,
		UserId:   userObjectId,
		Reaction: reaction.Reaction,
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": messageObjectId}, bson.M{"$push": bson.M{"reactions": mongoReaction}})
	if err != nil {
		return err
	}

	return nil
}
