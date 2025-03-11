package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	databaseName   = "supchat"
	collectionName = "chat_direct"
)

type MongoChatDirect struct {
	Id        bson.ObjectID              `bson:"_id"`
	SenderId  bson.ObjectID              `bson:"sender_id"`
	User1Id   bson.ObjectID              `bson:"user1Id"`
	User2Id   bson.ObjectID              `bson:"user2Id"`
	Content   string                     `bson:"content"`
	Reactions []*MongoChatDirectReaction `bson:"reactions"`
	CreatedAt time.Time                  `bson:"created_at"`
	UpdatedAt time.Time                  `bson:"updated_at"`
}

type MongoChatDirectReaction struct {
	Id       bson.ObjectID   `bson:"_id"`
	Users    []bson.ObjectID `bson:"user_ids"`
	Reaction string          `bson:"reaction"`
}

type MongoChatDirectRepositoryDeps struct {
	uberdig.In
	Mapper mapper.Mapper[*MongoChatDirect, *entity.ChatDirect]
	Client *mongo.Client
}

type MongoChatDirectRepository struct {
	deps MongoChatDirectRepositoryDeps
}

func NewMongoChatDirectRepository(deps MongoChatDirectRepositoryDeps) ChatDirectRepository {
	return &MongoChatDirectRepository{deps: deps}
}

func (m MongoChatDirectRepository) Create(ctx context.Context, chatDirect *entity.ChatDirect) error {
	if chatDirect.Id == "" {
		chatDirect.Id = entity.ChatDirectId(bson.NewObjectID().Hex())
	}
	if chatDirect.CreatedAt.IsZero() {
		now := time.Now()
		chatDirect.CreatedAt = now
		chatDirect.UpdatedAt = now
	}

	mongoChatDirect, err := m.deps.Mapper.MapFromEntity(chatDirect)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoChatDirect)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoChatDirectRepository) ListRecentChats(ctx context.Context) ([]*entity.ChatDirect, error) {
	pipeline := mongo2.Pipeline{
		// Étape 1 : Ajouter un champ `pairKey` qui normalise l'ordre des IDs
		{{"$addFields", bson.D{
			{"sortedIds", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{{"$lt", bson.A{"$user1Id", "$user2Id"}}}},
					{"then", bson.A{"$user1Id", "$user2Id"}},
					{"else", bson.A{"$user2Id", "$user1Id"}},
				}},
			}},
		}}},
		// Étape 2 : Grouper par `sortedIds` en prenant le message le plus récent
		{{"$group", bson.D{
			{"_id", "$sortedIds"},
			{"latestMessage", bson.D{{"$last", "$$ROOT"}}},
		}}},
		// Étape 3 : Trier par `updated_at` du dernier message
		{{"$sort", bson.D{{"latestMessage.updated_at", -1}}}},
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var groupedMongoChatDirects []struct {
		LatestMessage MongoChatDirect `bson:"latestMessage"`
	}

	if err = cursor.All(ctx, &groupedMongoChatDirects); err != nil {
		return nil, err
	}

	chatDirects := make([]*entity.ChatDirect, len(groupedMongoChatDirects))
	for i, groupedMongoChatDirect := range groupedMongoChatDirects {
		chatDirect, err := m.deps.Mapper.MapToEntity(&groupedMongoChatDirect.LatestMessage)
		if err != nil {
			return nil, err
		}
		chatDirects[i] = chatDirect
	}

	return chatDirects, nil
}

func (m MongoChatDirectRepository) ListByUser(ctx context.Context, user1Id, user2Id user_entity.UserId) ([]*entity.ChatDirect, error) {
	user1IdHex, err := bson.ObjectIDFromHex(string(user1Id))
	if err != nil {
		return nil, err
	}

	user2IdHex, err := bson.ObjectIDFromHex(string(user2Id))
	if err != nil {
		return nil, err
	}

	// Filter must match both ways
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"user1Id", user1IdHex}, {"user2Id", user2IdHex}},
			bson.D{{"user1Id", user2IdHex}, {"user2Id", user1IdHex}},
		},
		},
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var mongoChatDirects []*MongoChatDirect
	if err = cursor.All(ctx, &mongoChatDirects); err != nil {
		return nil, err
	}

	chatDirects := make([]*entity.ChatDirect, len(mongoChatDirects))
	for i, mongoChatDirect := range mongoChatDirects {
		chatDirect, err := m.deps.Mapper.MapToEntity(mongoChatDirect)
		if err != nil {
			return nil, err
		}
		chatDirects[i] = chatDirect
	}

	return chatDirects, nil
}
