package export_data_user_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type ExportDataUserWorkspaceUseCase struct {
	workspaceRepository repository.WorkspaceRepository
}

func NewExportDataUserWorkspaceUseCase(
	workspaceRepository repository.WorkspaceRepository,
) *ExportDataUserWorkspaceUseCase {
	return &ExportDataUserWorkspaceUseCase{workspaceRepository: workspaceRepository}
}

func (u *ExportDataUserWorkspaceUseCase) Execute(ctx context.Context, userId entity.UserId) ([]ExportableWorkspaceData, error) {
	workspaces, err := u.workspaceRepository.ListAllWorkspacesByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	exported := make([]ExportableWorkspaceData, len(workspaces))
	for i, workspace := range workspaces {
		exported[i] = ExportableWorkspaceData{
			Id:      workspace.Id,
			Name:    workspace.Name,
			Topic:   workspace.Topic,
			Type:    workspace.Type,
			OwnerId: workspace.OwnerId,
			TrucID:  workspace.TrucID,
		}
	}
	return exported, nil
}

type ExportableWorkspaceData struct {
	Id      entity2.WorkspaceId   `json:"id"`
	Name    string                `json:"name"`
	Topic   string                `json:"topic,omitempty"`
	Type    entity2.WorkspaceType `json:"type"`
	OwnerId entity.UserId         `json:"ownerId"`
	TrucID  entity.UserId         `json:"trucId,omitempty"`
}
