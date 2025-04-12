package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	databaseName   = "supchat"
	collectionName = "user_status"
)

var (
	UnknownUserStatusErr = errors.New("unknown user status")
)

type MongoUserStatusRepository struct {
	client *mongo.Client
}

type MongoUserStatus struct {
	UserId     bson.ObjectID `bson:"user_id"`
	UserStatus string        `bson:"user_status"`
}

func NewMongoUserStatusRepository(client *mongo.Client) UserStatusRepository {
	return &MongoUserStatusRepository{client: client}
}

func (m MongoUserStatusRepository) Get(ctx context.Context, userId user_entity.UserId) (*entity.UserStatus, error) {
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	var mongoUserStatus MongoUserStatus
	err = m.client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"user_id": userObjectId}).Decode(&mongoUserStatus)
	if err != nil {
		return nil, err
	}

	parsedUserStatus := entity.ParseStatus(mongoUserStatus.UserStatus)
	if parsedUserStatus == entity.StatusUnknown {
		return nil, UnknownUserStatusErr
	}

	return &entity.UserStatus{
		UserId: user_entity.UserId(mongoUserStatus.UserId.Hex()),
		Status: parsedUserStatus,
	}, nil
}

func (m MongoUserStatusRepository) Save(ctx context.Context, userStatus *entity.UserStatus) error {
	userObjectId, err := bson.ObjectIDFromHex(userStatus.UserId.String())
	if err != nil {
		return err
	}

	mongoUserStatus := MongoUserStatus{
		UserId:     userObjectId,
		UserStatus: userStatus.Status.String(),
	}
	_, err = m.client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"user_id": userObjectId}, bson.M{
		"$set": mongoUserStatus,
	}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
