package mongo

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity" // Assurez-vous d'importer le bon chemin pour UserId
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	databaseName   = "supchat"
	collectionName = "notifications"
)

type MongoNotificationRepositoryDeps struct {
	uberdig.In
	Client             *mongo.Client
	NotificationMapper mapper.Mapper[*MongoNotification, *entity.Notification]
}

type MongoNotificationRepository struct {
	deps MongoNotificationRepositoryDeps
}

type MongoNotification struct {
	Id        bson.ObjectID `bson:"_id"`
	UserId    string        `bson:"user_id"`
	Content   string        `bson:"content"`
	IsRead    bool          `bson:"is_read"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func NewMongoNotificationRepository(deps MongoNotificationRepositoryDeps) repository.NotificationRepository {
	return &MongoNotificationRepository{deps: deps}
}

func (r MongoNotificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	notification.Id = entity.NotificationId(bson.NewObjectID().Hex())
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = notification.CreatedAt
	mongoEntity, err := r.deps.NotificationMapper.MapFromEntity(notification)
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoEntity)
	if err != nil {
		return err
	}

	return nil
}

func (r MongoNotificationRepository) GetById(ctx context.Context, notificationId entity.NotificationId) (notification *entity.Notification, err error) {
	notificationObjectId, err := bson.ObjectIDFromHex(notificationId.String())
	if err != nil {
		return nil, err
	}

	var mongoNotification *MongoNotification
	err = r.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": notificationObjectId}).Decode(&mongoNotification)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, repository.NotificationNotFoundErr
		}
		return nil, err
	}

	notification, err = r.deps.NotificationMapper.MapToEntity(mongoNotification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (r MongoNotificationRepository) List(ctx context.Context, userId user_entity.UserId) (notifications []*entity.Notification, err error) {
	cursor, err := r.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{"user_id": userId.String()})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var mongoNotification *MongoNotification
		err = cursor.Decode(&mongoNotification)
		if err != nil {
			return nil, err
		}

		notification, err := r.deps.NotificationMapper.MapToEntity(mongoNotification)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r MongoNotificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	mongoEntity, err := r.deps.NotificationMapper.MapFromEntity(notification)
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"_id": mongoEntity.Id}, bson.M{"$set": mongoEntity})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return repository.NotificationNotFoundErr
		}
		return err
	}

	return nil
}

func (r MongoNotificationRepository) Delete(ctx context.Context, notificationId entity.NotificationId) error {
	objectId, err := bson.ObjectIDFromHex(notificationId.String())
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).Collection(collectionName).DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return repository.NotificationNotFoundErr
		}
		return err
	}

	return nil
}
