package create_channel

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workpace_members"
	uberdig "go.uber.org/dig"
)

type NotifyWebSocketObserverDeps struct {
	uberdig.In
	ListWorkspaceMembersUseCase *list_workpace_members.ListWorkspaceMembersUseCase
	WsServer                    *websocket.WsServer
	Logger                      logger.Logger
}

type NotifyWebSocketObserver struct {
	deps NotifyWebSocketObserverDeps
}

func NewNotifyWebSocketObserver(deps NotifyWebSocketObserverDeps) CreateChannelObserver {
	return &NotifyWebSocketObserver{deps: deps}
}

func (o *NotifyWebSocketObserver) ChannelCreated(channel *channel_entity.Channel) {
	// First, get all the clients that are in the workspace and connected to the websocket.

	workspaceMembers, err := o.deps.ListWorkspaceMembersUseCase.Execute(context.Background(), channel.WorkspaceId)
	if err != nil {
		o.deps.Logger.Error().Err(err).Msg("Error on getting workspace members")
		return
	}

	for _, member := range workspaceMembers {
		o.deps.WsServer.IterateClients(func(client *websocket.Client) (stop bool) {
			// Then, notify all the clients that a new channel has been created.
			if client.UserId == member.UserId {
				client.SendMessage(
					websocket.NewMessageBuilder().
						WithAction(websocket.ChannelCreatedAction).
						WithPayload(ChannelPayload{
							Id:    string(channel.Id),
							Name:  channel.Name,
							Topic: channel.Topic,
						}).
						Build(),
				)
			}

			return
		})
	}

}

type ChannelPayload struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
