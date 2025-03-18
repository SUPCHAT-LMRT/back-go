package count_messages_by_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type CountMessagesByWorkspaceUseCase struct {
	repository repository.ChannelMessageRepository
}

func NewCountMessagesUseCase(repository repository.ChannelMessageRepository) *CountMessagesByWorkspaceUseCase {
	return &CountMessagesByWorkspaceUseCase{repository: repository}
}

func (u CountMessagesByWorkspaceUseCase) Execute(ctx context.Context, workspaceId workspace_entity.WorkspaceId) (uint, error) {
	return u.repository.CountByWorkspace(ctx, workspaceId)
}
