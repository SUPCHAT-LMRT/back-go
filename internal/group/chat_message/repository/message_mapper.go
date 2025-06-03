package repository

import (
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GroupChatMessageMapper struct{}

func NewGroupChatMessageMapper() mapper.Mapper[*MongoGroupChatMessage, *entity.GroupChatMessage] {
	return &GroupChatMessageMapper{}
}

func (m GroupChatMessageMapper) MapToEntity(
	mongo *MongoGroupChatMessage,
) (*entity.GroupChatMessage, error) {
	return &entity.GroupChatMessage{
		Id:        entity.GroupChatMessageId(mongo.Id.Hex()),
		GroupId:   group_entity.GroupId(mongo.GroupId.Hex()),
		AuthorId:  user_entity.UserId(mongo.AuthorId.Hex()),
		Content:   mongo.Content,
		CreatedAt: mongo.CreatedAt,
	}, nil
}

func (m GroupChatMessageMapper) MapFromEntity(
	entity *entity.GroupChatMessage,
) (*MongoGroupChatMessage, error) {
	messageObjectId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		return nil, err
	}
	groupObjectId, err := bson.ObjectIDFromHex(string(entity.GroupId))
	if err != nil {
		return nil, err
	}
	authorObjectId, err := bson.ObjectIDFromHex(string(entity.AuthorId))
	if err != nil {
		return nil, err
	}

	return &MongoGroupChatMessage{
		Id:        messageObjectId,
		GroupId:   groupObjectId,
		AuthorId:  authorObjectId,
		Content:   entity.Content,
		CreatedAt: entity.CreatedAt,
	}, nil
}
