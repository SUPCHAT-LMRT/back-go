package websocket

import (
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type WebSocketDeps struct {
	uberdig.In
	GetChannelUseCase          *get_channel.GetChannelUseCase
	GetWorkspaceMemberUseCase  *get_workpace_member.GetWorkspaceMemberUseCase
	ListChannelMessagesUseCase *list_messages.ListChannelMessagesUseCase
	GetUserByIdUseCase         *get_by_id.GetUserByIdUseCase
	SendMessageObservers       []SendMessageObserver `group:"send_message_observers"`
	ToggleReactionUseCase      *toggle_reaction.ToggleReactionUseCase
	Logger                     logger.Logger
	RedisClient                *redis.Client
}
