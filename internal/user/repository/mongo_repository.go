package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	databaseName   = "supchat"
	collectionName = "users"
)

type MongoUserRepositoryDeps struct {
	uberdig.In
	Client     *mongo.Client
	UserMapper mapper.Mapper[*MongoUser, *entity.User]
}

type MongoUserRepository struct {
	deps MongoUserRepositoryDeps
}

type MongoUser struct {
	Id        bson.ObjectID `bson:"_id"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	CreatedAt time.Time     `bson:"created_at"`
}

func NewMongoUserRepository(deps MongoUserRepositoryDeps) UserRepository {
	return &MongoUserRepository{deps: deps}
}

func (m MongoUserRepository) Create(ctx context.Context, user *entity.User) error {
	user.Id = entity.UserId(bson.NewObjectID().Hex())
	user.CreatedAt = time.Now()

	mongoEntity, err := m.deps.UserMapper.MapFromEntity(user)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoEntity)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoUserRepository) GetById(ctx context.Context, userId entity.UserId) (user *entity.User, err error) {
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	var mongoUser *MongoUser
	err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": userObjectId}).Decode(&mongoUser)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, UserNotFoundErr
		}
		return nil, err
	}

	user, err = m.deps.UserMapper.MapToEntity(mongoUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m MongoUserRepository) GetByEmail(ctx context.Context, userEmail string, options ...GetUserOptionFunc) (user *entity.User, err error) {
	var mongoUser *MongoUser
	err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"email": userEmail}).Decode(&mongoUser)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, UserNotFoundErr
		}
		return nil, err
	}

	user, err = m.deps.UserMapper.MapToEntity(mongoUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m MongoUserRepository) List(ctx context.Context) (users []*entity.User, err error) {
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var mongoUser *MongoUser
		err = cursor.Decode(&mongoUser)
		if err != nil {
			return nil, err
		}

		user, err := m.deps.UserMapper.MapToEntity(mongoUser)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (m MongoUserRepository) Update(ctx context.Context, user *entity.User) error {
	mongoEntity, err := m.deps.UserMapper.MapFromEntity(user)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"_id": mongoEntity.Id}, bson.M{"$set": mongoEntity})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return UserNotFoundErr
		}
		return err
	}

	return nil
}

func (m MongoUserRepository) Delete(ctx context.Context, userId entity.UserId) error {
	_, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return UserNotFoundErr
		}
		return err
	}

	return nil
}
