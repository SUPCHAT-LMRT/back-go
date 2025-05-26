package repository

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "oauth_connections"
)

type MongoOauthConnectionRepositoryDeps struct {
	uberdig.In
	Client *mongo.Client
	Mapper mapper.Mapper[*MongoOauthConnection, *entity.OauthConnection]
}
type MongoOauthConnectionRepository struct {
	deps MongoOauthConnectionRepositoryDeps
}

type MongoOauthConnection struct {
	Id          bson.ObjectID `bson:"_id"`
	UserId      bson.ObjectID `bson:"user_id"`
	Provider    string        `bson:"provider"`
	OauthEmail  string        `bson:"oauth_email"`
	OauthUserId string        `bson:"oauth_user_id"`
}

func NewMongoOauthConnectionRepository(
	deps MongoOauthConnectionRepositoryDeps,
) OauthConnectionRepository {
	return &MongoOauthConnectionRepository{deps: deps}
}

func (m MongoOauthConnectionRepository) CreateOauthConnection(
	ctx context.Context,
	connection *entity.OauthConnection,
) error {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)

	connection.Id = entity.OauthConnectionId(bson.NewObjectID().Hex())
	mongoConnection, err := m.deps.Mapper.MapFromEntity(connection)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, mongoConnection)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoOauthConnectionRepository) GetOauthConnectionByUserId(
	ctx context.Context,
	userId string,
) (*entity.OauthConnection, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	filter := bson.M{"oauth_user_id": userId}
	var mongoConnection MongoOauthConnection
	err := collection.FindOne(ctx, filter).Decode(&mongoConnection)
	if err != nil {
		return nil, err
	}

	connection, err := m.deps.Mapper.MapToEntity(&mongoConnection)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (m MongoOauthConnectionRepository) ListOauthConnectionsByUser(
	ctx context.Context,
	userId user_entity.UserId,
) ([]*entity.OauthConnection, error) {
	collection := m.deps.Client.Client.Database(databaseName).Collection(collectionName)
	filter := bson.M{"user_id": userId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var connections []*entity.OauthConnection
	for cursor.Next(ctx) {
		var mongoConnection MongoOauthConnection
		if err = cursor.Decode(&mongoConnection); err != nil {
			return nil, err
		}

		connection, err := m.deps.Mapper.MapToEntity(&mongoConnection)
		if err != nil {
			return nil, err
		}

		connections = append(connections, connection)
	}

	return connections, nil
}
