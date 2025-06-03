package repository

import (
	"context"
	"time"

	time_series_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type MessageSentTimeSeriesWorkspaceRepository interface {
	Create(
		ctx context.Context,
		sentAt time.Time,
		metadata time_series_entity.MessageSentMetadata,
	) error
	GetMinutelyByWorkspace(
		ctx context.Context,
		workspaceId entity.WorkspaceId,
		from, to time.Time,
	) ([]*time_series_entity.MessageSent, error)
}
