package get_minutely

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/repository"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetMinutelyMessageSentUseCase struct {
	repository repository.MessageSentTimeSeriesWorkspaceRepository
}

func NewGetMinutelyMessageSentUseCase(
	messageSentTimeSeriesWorkspaceRepository repository.MessageSentTimeSeriesWorkspaceRepository,
) *GetMinutelyMessageSentUseCase {
	return &GetMinutelyMessageSentUseCase{repository: messageSentTimeSeriesWorkspaceRepository}
}

func (u *GetMinutelyMessageSentUseCase) Execute(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
	from, to time.Time,
) ([]*entity.MessageSent, error) {
	return u.repository.GetMinutelyByWorkspace(ctx, workspaceId, from, to)
}
