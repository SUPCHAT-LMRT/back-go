package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatDirectMapper struct{}

func NewChatDirectMapper() mapper.Mapper[*MongoChatDirect, *chat_direct_entity.ChatDirect] {
	return &ChatDirectMapper{}
}

func (m ChatDirectMapper) MapToEntity(mongo *MongoChatDirect) (*chat_direct_entity.ChatDirect, error) {
	return &chat_direct_entity.ChatDirect{
		Id:        chat_direct_entity.ChatDirectId(mongo.Id.Hex()),
		User1Id:   entity.UserId(mongo.User1Id.Hex()),
		User2Id:   entity.UserId(mongo.User2Id.Hex()),
		CreatedAt: mongo.CreatedAt,
		UpdatedAt: mongo.UpdatedAt,
	}, nil
}

func (m ChatDirectMapper) MapFromEntity(entity *chat_direct_entity.ChatDirect) (*MongoChatDirect, error) {
	chatObjectId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		return nil, err
	}

	user1Id, err := bson.ObjectIDFromHex(string(entity.User1Id))
	if err != nil {
		return nil, err
	}

	user2Id, err := bson.ObjectIDFromHex(string(entity.User2Id))
	if err != nil {
		return nil, err
	}

	return &MongoChatDirect{
		Id:        chatObjectId,
		User1Id:   user1Id,
		User2Id:   user2Id,
		UpdatedAt: entity.UpdatedAt,
		CreatedAt: entity.CreatedAt,
	}, nil
}
