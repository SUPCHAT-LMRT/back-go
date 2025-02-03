package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	databaseName   = "supchat"
	collectionName = "chat_direct"
)

type MongoChatDirect struct {
	Id        bson.ObjectID `bson:"_id"`
	User1Id   bson.ObjectID `bson:"user1Id"`
	User2Id   bson.ObjectID `bson:"user2Id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
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
	chatDirect.Id = entity.ChatDirectId(bson.NewObjectID().Hex())
	chatDirect.CreatedAt = time.Now()
	chatDirect.UpdatedAt = chatDirect.CreatedAt

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

func (m MongoChatDirectRepository) ListRecentGroups(ctx context.Context) ([]*entity.ChatDirect, error) {
	opts := options.Find().SetSort(bson.D{{"updated_at", 1}})
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.D{}, opts)
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
