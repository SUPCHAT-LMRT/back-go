package websocket

import (
	"github.com/supchat-lmrt/back-go/internal/back_identifier/usecase"
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/redis"
	toggle_direct_message_reaction "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	toggle_channel_message_reaction "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	uberdig "go.uber.org/dig"
)

type WebSocketDeps struct {
	uberdig.In
	GetChannelUseCase                   *get_channel.GetChannelUseCase
	GetWorkspaceMemberUseCase           *get_workpace_member.GetWorkspaceMemberUseCase
	ListChannelMessagesUseCase          *list_messages.ListChannelMessagesUseCase
	GetUserByIdUseCase                  *get_by_id.GetUserByIdUseCase
	SendChannelMessageObservers         []SendChannelMessageObserver `group:"send_channel_message_observers"`
	SendDirectMessageObservers          []SendDirectMessageObserver  `group:"send_direct_message_observers"`
	ToggleReactionChannelMessageUseCase *toggle_channel_message_reaction.ToggleReactionChannelMessageUseCase
	ToggleReactionDirectMessageUseCase  *toggle_direct_message_reaction.ToggleReactionDirectMessageUseCase
	Logger                              logger.Logger
	RedisClient                         *redis.Client
	GetBackIdentifierUseCase            *usecase.GetBackIdentifierUseCase
	EventBus                            *event.EventBus
}
