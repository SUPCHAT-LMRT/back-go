package mongo

import (
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrMissingDirectMessageData = errors.New("missing direct message data in notification")
var ErrMissingChannelMessageData = errors.New("missing direct message data in notification")

type MongoNotificationMapper struct{}

func NewMongoNotificationMapper() mapper.Mapper[*MongoNotification, *entity.Notification] {
	return &MongoNotificationMapper{}
}

func (m MongoNotificationMapper) MapFromEntity(notification *entity.Notification) (*MongoNotification, error) {
	notificationObjectId, err := bson.ObjectIDFromHex(string(notification.Id))
	if err != nil {
		return nil, err
	}
	notificationUserObjectId, err := bson.ObjectIDFromHex(string(notification.UserId))
	if err != nil {
		return nil, err
	}

	mongonotification := MongoNotification{
		Id:                  notificationObjectId,
		UserId:              notificationUserObjectId,
		Type:                notification.Type.String(),
		IsRead:              notification.IsRead,
		CreatedAt:           notification.CreatedAt,
		DirectMessageData:   nil,
		ChannelMessageData:  nil,
		WorkspaceInviteData: nil,
	}
	switch notification.Type {
	case entity.NotificationTypeDirectMessage:
		if notification.DirectMessageData == nil {
			return nil, ErrMissingDirectMessageData
		}
		notificationSenderObjectId, err := bson.ObjectIDFromHex(string(notification.DirectMessageData.SenderId))
		if err != nil {
			return nil, err
		}
		notificationmessageObjectId, err := bson.ObjectIDFromHex(string(notification.DirectMessageData.MessageId))
		if err != nil {
			return nil, err
		}
		mongonotification.DirectMessageData = &DirectMessageNotificationData{
			SenderId:  notificationSenderObjectId,
			MessageId: notificationmessageObjectId,
		}
	case entity.NotificationTypeChannelMessage:
		if notification.ChannelMessageData == nil {
			return nil, ErrMissingChannelMessageData
		}
		notificationSenderObjectId, err := bson.ObjectIDFromHex(string(notification.ChannelMessageData.SenderId))
		if err != nil {
			return nil, err
		}
		notificationmessageObjectId, err := bson.ObjectIDFromHex(string(notification.ChannelMessageData.MessageId))
		if err != nil {
			return nil, err
		}
		notificationChannelObjectId, err := bson.ObjectIDFromHex(string(notification.ChannelMessageData.ChannelId))
		if err != nil {
			return nil, err
		}
		notificationWorkspaceObjectId, err := bson.ObjectIDFromHex(string(notification.ChannelMessageData.WorkspaceId))
		if err != nil {
			return nil, err
		}
		mongonotification.ChannelMessageData = &ChannelMessageNotificationData{
			SenderId: notificationSenderObjectId,
			//SenderAvatarUrl: notification.ChannelMessageData.SenderAvatarUrl,
			ChannelId:   notificationChannelObjectId,
			WorkspaceId: notificationWorkspaceObjectId,
			MessageId:   notificationmessageObjectId,
		}
	case entity.NotificationTypeWorkspaceInvite:
		if notification.ChannelMessageData == nil {
			return nil, ErrMissingChannelMessageData
		}
		notificationInviterObjectId, err := bson.ObjectIDFromHex(string(notification.WorkspaceInviteData.InviterId))
		if err != nil {
			return nil, err
		}
		notificationWorkspaceObjectId, err := bson.ObjectIDFromHex(string(notification.WorkspaceInviteData.WorkspaceId))
		if err != nil {
			return nil, err
		}
		mongonotification.WorkspaceInviteData = &WorkspaceInviteNotificationData{
			InviterId:   notificationInviterObjectId,
			WorkspaceId: notificationWorkspaceObjectId,
		}

	}
	return &mongonotification, nil
}

func (m MongoNotificationMapper) MapToEntity(databaseNotification *MongoNotification) (*entity.Notification, error) {
	notification := &entity.Notification{
		Id:        entity.NotificationId(databaseNotification.Id.Hex()),
		UserId:    user_entity.UserId(databaseNotification.UserId.Hex()),
		Type:      entity.NotificationType(databaseNotification.Type),
		IsRead:    databaseNotification.IsRead,
		CreatedAt: databaseNotification.CreatedAt,
	}

	switch notification.Type {
	case entity.NotificationTypeDirectMessage:
		if databaseNotification.DirectMessageData == nil {
			return nil, ErrMissingDirectMessageData
		}
		notification.DirectMessageData = &entity.DirectMessageNotificationData{
			SenderId:  user_entity.UserId(databaseNotification.DirectMessageData.SenderId.Hex()),
			MessageId: chat_direct_entity.ChatDirectId(databaseNotification.DirectMessageData.MessageId.Hex()),
		}
	case entity.NotificationTypeChannelMessage:
		if databaseNotification.ChannelMessageData == nil {
			return nil, ErrMissingChannelMessageData
		}
		notification.ChannelMessageData = &entity.ChannelMessageNotificationData{
			SenderId: user_entity.UserId(databaseNotification.ChannelMessageData.SenderId.Hex()),
			//SenderAvatarUrl: databaseNotification.ChannelMessageData.SenderAvatarUrl,
			ChannelId:   channel_entity.ChannelId(databaseNotification.ChannelMessageData.ChannelId.Hex()),
			WorkspaceId: workspace_entity.WorkspaceId(databaseNotification.ChannelMessageData.WorkspaceId.Hex()),
			MessageId:   channel_message_entity.ChannelMessageId(databaseNotification.ChannelMessageData.MessageId.Hex()),
		}
	case entity.NotificationTypeWorkspaceInvite:
		if databaseNotification.WorkspaceInviteData == nil {
			return nil, ErrMissingChannelMessageData
		}
		notification.WorkspaceInviteData = &entity.WorkspaceInviteNotificationData{
			InviterId:   user_entity.UserId(databaseNotification.WorkspaceInviteData.InviterId.Hex()),
			WorkspaceId: workspace_entity.WorkspaceId(databaseNotification.WorkspaceInviteData.WorkspaceId.Hex()),
		}
	}
	return notification, nil
}
