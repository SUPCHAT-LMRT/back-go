package repository

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "group_chat_messages"
)

type MongoGroupChatMessage struct {
	Id        bson.ObjectID           `bson:"_id"`
	GroupId   bson.ObjectID           `bson:"group_id"`
	SenderId  bson.ObjectID           `bson:"sender_id"`
	Content   string                  `bson:"content"`
	Reactions []*MongoMessageReaction `bson:"reactions"`
	CreatedAt time.Time               `bson:"created_at"`
	UpdatedAt time.Time               `bson:"updated_at"`
}

type MongoMessageReaction struct {
	Id       bson.ObjectID   `bson:"_id"`
	Users    []bson.ObjectID `bson:"users"`
	Reaction string          `bson:"reaction"`
}

type MongoGroupChatRepositoryDeps struct {
	uberdig.In
	Mapper mapper.Mapper[*MongoGroupChatMessage, *entity.GroupChatMessage]
	Client *mongo.Client
}

type MongoGroupChatRepository struct {
	deps MongoGroupChatRepositoryDeps
}

func NewMongoGroupChatRepository(deps MongoGroupChatRepositoryDeps) ChatMessageRepository {
	return &MongoGroupChatRepository{deps: deps}
}

func (m MongoGroupChatRepository) Create(
	ctx context.Context,
	chatMessage *entity.GroupChatMessage,
) error {
	if chatMessage.Id == "" {
		chatMessage.Id = entity.GroupChatMessageId(bson.NewObjectID().Hex())
	}
	now := time.Now()
	if chatMessage.CreatedAt.IsZero() {
		chatMessage.CreatedAt = now
		chatMessage.UpdatedAt = now
	}
	if chatMessage.UpdatedAt.IsZero() {
		chatMessage.UpdatedAt = now
	}

	mongoMessage, err := m.deps.Mapper.MapFromEntity(chatMessage)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoMessage)

	return err
}

func (m MongoGroupChatRepository) GetLastMessage(
	ctx context.Context,
	groupId group_entity.GroupId,
) (*entity.GroupChatMessage, error) {
	groupObjId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return nil, err
	}

	findOptions := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(1)

	cursor, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.D{{"group_id", groupObjId}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoMessage MongoGroupChatMessage
	if !cursor.Next(ctx) {
		return nil, mongo2.ErrNoDocuments
	}
	if err := cursor.Decode(&mongoMessage); err != nil {
		return nil, err
	}

	return m.deps.Mapper.MapToEntity(&mongoMessage)
}

func (m MongoGroupChatRepository) IsFirstMessage(
	ctx context.Context,
	groupId group_entity.GroupId,
) (bool, error) {
	groupObjId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return false, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.D{{"group_id", groupObjId}})
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (m MongoGroupChatRepository) ListMessages(
	ctx context.Context,
	groupId group_entity.GroupId,
	params ListMessagesQueryParams,
) ([]*entity.GroupChatMessage, error) {
	groupObjId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"group_id", groupObjId}}

	if params.AroundMessageId != "" {
		// Handle around message logic
		messageId, err := bson.ObjectIDFromHex(params.AroundMessageId.String())
		if err != nil {
			return nil, err
		}

		// Get the message to find its timestamp
		var message MongoGroupChatMessage
		err = m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			FindOne(ctx, bson.D{{"_id", messageId}}).
			Decode(&message)
		if err != nil {
			return nil, err
		}

		// Fetch messages around the target message's timestamp
		halfLimit := params.Limit / 2

		beforeFilter := append(
			filter,
			bson.E{Key: "created_at", Value: bson.M{"$lt": message.CreatedAt}},
		)

		// Get messages before the target
		beforeCursor, err := m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			Find(ctx, beforeFilter, options.Find().
				SetSort(bson.D{{"created_at", -1}}).
				SetLimit(int64(halfLimit)))
		if err != nil {
			return nil, err
		}
		defer beforeCursor.Close(ctx)

		afterFilter := append(
			filter,
			bson.E{Key: "created_at", Value: bson.M{"gt": message.CreatedAt}},
		)

		// Get messages after the target
		afterCursor, err := m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			Find(ctx, afterFilter, options.Find().
				SetSort(bson.D{{"created_at", 1}}).
				SetLimit(int64(halfLimit)))
		if err != nil {
			return nil, err
		}
		defer afterCursor.Close(ctx)

		var results []*entity.GroupChatMessage

		// Process messages before
		var beforeMessages []*entity.GroupChatMessage
		for beforeCursor.Next(ctx) {
			var mongoMsg MongoGroupChatMessage
			if err := beforeCursor.Decode(&mongoMsg); err != nil {
				return nil, err
			}
			entityMsg, err := m.deps.Mapper.MapToEntity(&mongoMsg)
			if err != nil {
				return nil, err
			}
			beforeMessages = append(beforeMessages, entityMsg)
		}

		// Reverse before messages to maintain chronological order
		for i := len(beforeMessages) - 1; i >= 0; i-- {
			results = append(results, beforeMessages[i])
		}

		// Process messages after
		for afterCursor.Next(ctx) {
			var mongoMsg MongoGroupChatMessage
			if err := afterCursor.Decode(&mongoMsg); err != nil {
				return nil, err
			}
			entityMsg, err := m.deps.Mapper.MapToEntity(&mongoMsg)
			if err != nil {
				return nil, err
			}
			results = append(results, entityMsg)
		}

		return results, nil
	} else {
		// Standard pagination logic
		findOptions := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: -1}}).
			SetLimit(int64(params.Limit))

		if !params.Before.Equal(time.Time{}) {
			filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$lt": params.Before}})
		} else if !params.After.Equal(time.Time{}) {
			filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$gt": params.After}})
		}

		cursor, err := m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			Find(ctx, filter, findOptions)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var mongoGroupChatMessages []*MongoGroupChatMessage
		if err = cursor.All(ctx, &mongoGroupChatMessages); err != nil {
			return nil, err
		}

		results := make([]*entity.GroupChatMessage, len(mongoGroupChatMessages))
		for i, mongoGroupChatMessage := range mongoGroupChatMessages {
			entityMsg, err := m.deps.Mapper.MapToEntity(mongoGroupChatMessage)
			if err != nil {
				return nil, err
			}
			results[i] = entityMsg
		}

		return results, nil
	}
}

func (m MongoGroupChatRepository) ToggleReaction(
	ctx context.Context,
	messageId entity.GroupChatMessageId,
	userId user_entity.UserId,
	reaction string,
) (bool, error) {
	messageObjId, err := bson.ObjectIDFromHex(messageId.String())
	if err != nil {
		return false, err
	}

	userObjId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return false, err
	}

	// Check if reaction already exists
	filter := bson.D{
		{"_id", messageObjId},
		{"reactions.reaction", reaction},
	}

	// Try to add the user to the existing reaction
	update := bson.D{
		{"$addToSet", bson.D{
			{"reactions.$.users", userObjId},
		}},
	}

	result, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, filter, update)

	if err != nil {
		return false, err
	}

	// If user was added to the reaction
	if result.ModifiedCount > 0 {
		return true, nil
	}

	// Check if the user already reacted with this emoji
	filter = bson.D{
		{"_id", messageObjId},
		{"reactions.reaction", reaction},
		{"reactions.users", userObjId},
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, filter)

	if err != nil {
		return false, err
	}

	// If user already reacted, remove the reaction
	if count > 0 {
		update = bson.D{
			{"$pull", bson.D{
				{"reactions.$.users", userObjId},
			}},
		}

		_, err := m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			UpdateOne(ctx, bson.D{
				{"_id", messageObjId},
				{"reactions.reaction", reaction},
			}, update)

		if err != nil {
			return false, err
		}

		// Clean up empty reactions
		cleanup := bson.D{
			{"$pull", bson.D{
				{"reactions", bson.D{
					{"users", bson.D{{"$size", 0}}},
				}},
			}},
		}

		_, err = m.deps.Client.Client.Database(databaseName).
			Collection(collectionName).
			UpdateOne(ctx, bson.D{{"_id", messageObjId}}, cleanup)

		return false, err
	}

	// If reaction doesn't exist yet, create a new one
	reactionId := bson.NewObjectID()
	update = bson.D{
		{"$push", bson.D{
			{"reactions", bson.D{
				{"_id", reactionId},
				{"users", bson.A{userObjId}},
				{"reaction", reaction},
			}},
		}},
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.D{{"_id", messageObjId}}, update)

	return err == nil, err
}
