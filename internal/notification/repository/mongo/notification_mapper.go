package mongo

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoNotificationMapper struct{}

func NewMongoNotificationMapper() mapper.Mapper[*MongoNotification, *entity.Notification] {
	return &MongoNotificationMapper{}
}

func (m MongoNotificationMapper) MapFromEntity(
	entityNotification *entity.Notification,
) (*MongoNotification, error) {
	notificationObjectId, err := bson.ObjectIDFromHex(entityNotification.Id.String())
	if err != nil {
		return nil, err
	}

	return &MongoNotification{
		Id:        notificationObjectId,
		UserId:    entityNotification.UserId.String(),
		Content:   entityNotification.Content,
		IsRead:    entityNotification.IsRead,
		CreatedAt: entityNotification.CreatedAt,
		UpdatedAt: entityNotification.UpdatedAt,
	}, nil
}

func (m MongoNotificationMapper) MapToEntity(
	databaseNotification *MongoNotification,
) (*entity.Notification, error) {
	return &entity.Notification{
		Id:        entity.NotificationId(databaseNotification.Id.Hex()),
		UserId:    user_entity.UserId(databaseNotification.UserId),
		Content:   databaseNotification.Content,
		IsRead:    databaseNotification.IsRead,
		CreatedAt: databaseNotification.CreatedAt,
		UpdatedAt: databaseNotification.UpdatedAt,
	}, nil
}
