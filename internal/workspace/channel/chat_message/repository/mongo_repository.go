package repository

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
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
	UpdatedAt time.Time                      `bson:"updated_at"`
	Reactions []*MongoChannelMessageReaction `bson:"reactions"`
}

type MongoChannelMessageReaction struct {
	Id       bson.ObjectID   `bson:"_id"`
	Users    []bson.ObjectID `bson:"user_ids"`
	Reaction string          `bson:"reaction"`
}

func NewMongoChannelMessageRepository(
	deps MongoChannelMessageRepositoryDeps,
) ChannelMessageRepository {
	return &MongoChannelMessageRepository{deps: deps}
}

func (m MongoChannelMessageRepository) Create(
	ctx context.Context,
	message *entity.ChannelMessage,
) error {
	if message.Id == "" {
		message.Id = entity.ChannelMessageId(bson.NewObjectID().Hex())
	}
	now := time.Now()
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}
	if message.UpdatedAt.IsZero() {
		message.UpdatedAt = now
	}

	mongoMessage, err := m.deps.Mapper.MapFromEntity(message)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoMessage)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChannelMessageRepository) Get(
	ctx context.Context,
	id entity.ChannelMessageId,
) (*entity.ChannelMessage, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var mongoMessage MongoChannelMessage
	err = collection.FindOne(ctx, bson.M{"_id": messageObjectId}).Decode(&mongoMessage)
	if err != nil {
		return nil, err
	}

	message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
	if err != nil {
		return nil, err
	}

	return message, nil
}

//nolint:revive
func (m MongoChannelMessageRepository) ListByChannelId(
	ctx context.Context,
	channelId channel_entity.ChannelId,
	params ListByChannelIdQueryParams,
) ([]*entity.ChannelMessage, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	channelObjectId, err := bson.ObjectIDFromHex(string(channelId))
	if err != nil {
		return nil, err
	}

	var messages []*entity.ChannelMessage
	var aroundMessage MongoChannelMessage

	// Si AroundMessageId est défini, on récupère le message cible
	if params.AroundMessageId != "" {
		aroundObjectId, err := bson.ObjectIDFromHex(string(params.AroundMessageId))
		if err != nil {
			return nil, err
		}

		err = collection.FindOne(ctx, bson.M{"_id": aroundObjectId, "channel_id": channelObjectId}).
			Decode(&aroundMessage)
		if err != nil {
			return nil, err
		}

		// Récupération des messages avant le message cible
		beforeFilter := bson.M{
			"channel_id": channelObjectId,
			"created_at": bson.M{"$lt": aroundMessage.CreatedAt},
		}
		beforeOpts := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: -1}}).
			// Ordre décroissant pour prendre les plus récents en premier
			SetLimit(int64(params.Limit / 2)) // On récupère la moitié avant

		beforeCursor, err := collection.Find(ctx, beforeFilter, beforeOpts)
		if err != nil {
			return nil, err
		}
		defer beforeCursor.Close(ctx)

		var beforeMessages []*entity.ChannelMessage
		for beforeCursor.Next(ctx) {
			var mongoMessage MongoChannelMessage
			if err := beforeCursor.Decode(&mongoMessage); err != nil {
				return nil, err
			}
			message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
			if err != nil {
				return nil, err
			}
			beforeMessages = append(beforeMessages, message)
		}

		// Récupération des messages après le message cible
		afterFilter := bson.M{
			"channel_id": channelObjectId,
			"created_at": bson.M{"$gt": aroundMessage.CreatedAt},
		}
		afterOpts := options.Find().
			SetSort(bson.D{{Key: "created_at", Value: 1}}).
			// Ordre croissant pour prendre les plus anciens en premier
			SetLimit(int64(params.Limit / 2)) // On récupère la moitié après

		afterCursor, err := collection.Find(ctx, afterFilter, afterOpts)
		if err != nil {
			return nil, err
		}
		defer afterCursor.Close(ctx)

		var afterMessages []*entity.ChannelMessage
		for afterCursor.Next(ctx) {
			var mongoMessage MongoChannelMessage
			if err := afterCursor.Decode(&mongoMessage); err != nil {
				return nil, err
			}
			message, err := m.deps.Mapper.MapToEntity(&mongoMessage)
			if err != nil {
				return nil, err
			}
			afterMessages = append(afterMessages, message)
		}

		// On assemble le tout dans l'ordre chronologique
		for i := len(beforeMessages) - 1; i >= 0; i-- { // On remet l'ordre chronologique
			messages = append(messages, beforeMessages[i])
		}
		channelMessage, err := m.deps.Mapper.MapToEntity(&aroundMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, channelMessage) // Ajout du message cible
		messages = append(messages, afterMessages...)

		return messages, nil
	}

	// Sinon, on applique les filtres classiques
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // Trier par date décroissante
		SetLimit(int64(params.Limit))

	filter := bson.M{"channel_id": channelObjectId}

	if !params.Before.Equal(time.Time{}) {
		filter["created_at"] = bson.M{"$lt": params.Before}
	} else if !params.After.Equal(time.Time{}) {
		filter["created_at"] = bson.M{"$gt": params.After}
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mongoMessage MongoChannelMessage
		if err := cursor.Decode(&mongoMessage); err != nil {
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

// ToggleReaction toggles the reaction of a user to a message. (If the user has already reacted, it will remove the reaction, otherwise it will add the reaction.)
//
//nolint:revive
func (m MongoChannelMessageRepository) ToggleReaction(
	ctx context.Context,
	messageId entity.ChannelMessageId,
	userId user_entity.UserId,
	reaction string,
) (added bool, err error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(string(messageId))
	if err != nil {
		return false, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return false, err
	}

	var message MongoChannelMessage
	err = collection.FindOne(ctx, bson.M{"_id": messageObjectId}).Decode(&message)
	if err != nil {
		return false, err
	}

	updatedReactions := make([]*MongoChannelMessageReaction, 0)
	found := false
	removed := false

	for _, r := range message.Reactions {
		if r.Reaction == reaction {
			found = true
			updatedUsers := make([]bson.ObjectID, 0, len(r.Users))
			for _, uid := range r.Users {
				if uid.Hex() == userId.String() {
					removed = true
					continue
				}
				updatedUsers = append(updatedUsers, uid)
			}
			if len(updatedUsers) > 0 {
				r.Users = updatedUsers
				updatedReactions = append(updatedReactions, r)
			}
		} else {
			updatedReactions = append(updatedReactions, r)
		}
	}

	if !found {
		// Add new reaction if not found
		updatedReactions = append(updatedReactions, &MongoChannelMessageReaction{
			Id:       bson.NewObjectID(),
			Users:    []bson.ObjectID{userObjectId},
			Reaction: reaction,
		})
	} else if !removed {
		// Add user to existing reaction
		for _, r := range updatedReactions {
			if r.Reaction == reaction {
				r.Users = append(r.Users, userObjectId)
				break
			}
		}
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": messageObjectId},
		bson.M{"$set": bson.M{"reactions": updatedReactions}},
	)
	if err != nil {
		return false, err
	}

	return !removed, nil
}

func (m MongoChannelMessageRepository) CountByWorkspace(
	ctx context.Context,
	id workspace_entity.WorkspaceId,
) (uint, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	workspaceObjectId, err := bson.ObjectIDFromHex(string(id))
	if err != nil {
		return 0, err
	}

	lookupStage := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "workspaces_channels"},
				{Key: "localField", Value: "channel_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "channel_info"},
			},
		},
	}

	matchStage := bson.D{
		{
			Key: "$match", Value: bson.M{
				"channel_info.workspace_id": workspaceObjectId,
			},
		},
	}

	countStage := bson.D{
		{Key: "$count", Value: "total_messages"},
	}

	cursor, err := collection.Aggregate(ctx, mongo2.Pipeline{lookupStage, matchStage, countStage})
	if err != nil {
		return 0, err
	}

	var elementsCount struct {
		TotalMessages int `bson:"total_messages"`
	}
	if cursor.Next(ctx) {
		err = cursor.Decode(&elementsCount)
		if err != nil {
			return 0, err
		}
	}

	return uint(elementsCount.TotalMessages), nil
}

func (m *MongoChannelMessageRepository) ListAllMessagesByUser(ctx context.Context, userId user_entity.UserId) ([]*entity.ChannelMessage, error) {
	userIdHex, err := bson.ObjectIDFromHex(string(userId))
	if err != nil {
		return nil, err
	}

	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	filter := bson.M{
		"user_id": userIdHex,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoMessages []*MongoChannelMessage
	if err := cursor.All(ctx, &mongoMessages); err != nil {
		return nil, err
	}

	messages := make([]*entity.ChannelMessage, len(mongoMessages))
	for i, mongoMsg := range mongoMessages {
		msg, err := m.deps.Mapper.MapToEntity(mongoMsg)
		if err != nil {
			return nil, err
		}
		messages[i] = msg
	}
	return messages, nil
}

func (m MongoChannelMessageRepository) DeleteMessage(
	ctx context.Context,
	channelMessageId entity.ChannelMessageId,
) error {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(string(channelMessageId))
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": messageObjectId})
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChannelMessageRepository) UpdateMessage(
	ctx context.Context,
	msg *entity.ChannelMessage,
) error {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	messageObjectId, err := bson.ObjectIDFromHex(string(msg.Id))
	if err != nil {
		return err
	}

	mongoMessage, err := m.deps.Mapper.MapFromEntity(msg)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": messageObjectId},
		bson.M{"$set": mongoMessage},
	)
	if err != nil {
		return err
	}

	return nil
}
