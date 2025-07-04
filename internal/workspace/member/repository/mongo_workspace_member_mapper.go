package repository

import (
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoWorkspaceMemberMapper struct{}

func NewMongoWorkspaceMemberMapper() mapper.Mapper[*MongoWorkspaceMember, *entity.WorkspaceMember] {
	return &MongoWorkspaceMemberMapper{}
}

func (m MongoWorkspaceMemberMapper) MapFromEntity(
	workspaceMember *entity.WorkspaceMember,
) (*MongoWorkspaceMember, error) {
	memberObjectId, err := bson.ObjectIDFromHex(workspaceMember.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace member id to object id: %w", err)
	}

	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceMember.WorkspaceId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace id to object id: %w", err)
	}

	userObjectId, err := bson.ObjectIDFromHex(workspaceMember.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert user id to object id: %w", err)
	}

	return &MongoWorkspaceMember{
		Id:          memberObjectId,
		WorkspaceId: workspaceObjectId,
		UserId:      userObjectId,
	}, nil
}

func (m MongoWorkspaceMemberMapper) MapToEntity(
	databaseWorkspace *MongoWorkspaceMember,
) (*entity.WorkspaceMember, error) {
	return &entity.WorkspaceMember{
		Id:          entity.WorkspaceMemberId(databaseWorkspace.Id.Hex()),
		WorkspaceId: workspace_entity.WorkspaceId(databaseWorkspace.WorkspaceId.Hex()),
		UserId:      user_entity.UserId(databaseWorkspace.UserId.Hex()),
	}, nil
}
