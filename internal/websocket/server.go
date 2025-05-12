package websocket

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/event"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/outbound"
)

type WsServer struct {
	Deps           WebSocketDeps
	clients        map[*Client]bool
	Register       chan *Client
	Unregister     chan *Client
	rooms          map[*Room]bool
	backIdentifier string
}

func NewWsServer(deps WebSocketDeps) (*WsServer, error) {
	backIdentifier, err := deps.GetBackIdentifierUseCase.Execute(context.Background())
	if err != nil {
		return nil, err
	}

	server := &WsServer{
		Deps:           deps,
		clients:        make(map[*Client]bool),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		rooms:          make(map[*Room]bool),
		backIdentifier: backIdentifier,
	}

	server.Deps.EventBus.Subscribe(event.DirectChatMessageSavedEventType, func(evt event.Event) {
		messageSavedEvent, ok := evt.(*event.DirectChatMessageSavedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to DirectChatMessageSavedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("user1", messageSavedEvent.Message.User1Id.String()).
			Str("user2", messageSavedEvent.Message.User2Id.String()).Logger()

		user1Client := server.findClientByUserId(messageSavedEvent.Message.User1Id)
		if user1Client != nil {
			user2, err := deps.GetUserByIdUseCase.Execute(context.Background(), messageSavedEvent.Message.User2Id)
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to get user2")
				return
			}

			err = user1Client.SendMessage(&outbound.OutboundAddRecentDirectChat{
				OtherUserId: messageSavedEvent.Message.User2Id,
				ChatName:    user2.FullName(),
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send message to user1")
				return
			}
		}

		if messageSavedEvent.Message.User1Id == messageSavedEvent.Message.User2Id {
			// Skip sending message to self twice (user1 and user2 are the same)
			return
		}

		user2Client := server.findClientByUserId(messageSavedEvent.Message.User2Id)
		if user2Client != nil {
			user1, err := deps.GetUserByIdUseCase.Execute(context.Background(), messageSavedEvent.Message.User1Id)
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to get user1")
				return
			}

			err = user2Client.SendMessage(&outbound.OutboundAddRecentDirectChat{
				OtherUserId: messageSavedEvent.Message.User1Id,
				ChatName:    user1.FullName(),
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send message to user2")
				return
			}
		}
	})

	server.Deps.EventBus.Subscribe(event.UserStatusSavedEventType, func(evt event.Event) {
		userStatusSavedEvent, ok := evt.(*event.UserStatusSavedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to DirectChatMessageSavedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("user1", userStatusSavedEvent.UserStatus.UserId.String()).
			Str("status", userStatusSavedEvent.UserStatus.Status.String()).Logger()

		// Broadcast the status change to all clients
		server.IterateClients(func(client *Client) bool {
			// Check if the client is the one who changed the status, if so, send the status to self (a private status)
			// This is useful for the client to update its own status on different devices (changing status on phone and updating on web...)
			// Send to self, also the status change like for other users (to update the status on the recent chats for example) (this can be avoided by handling the self status in the frontend client)
			if client.UserId == userStatusSavedEvent.UserStatus.UserId {
				err = client.SendMessage(&outbound.OutboundSelfStatusUpdated{
					Status: userStatusSavedEvent.UserStatus.Status,
				})
				if err != nil {
					logg.Error().Err(err).
						Msg("failed to send status change to self client")
					return false
				}
			}

			err = client.SendMessage(&outbound.OutboundUserStatusUpdated{
				UserId: userStatusSavedEvent.UserStatus.UserId,
				Status: userStatusSavedEvent.UserStatus.Status.ToPublic(),
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send status change to client")
				return false
			}

			return true
		})
	})

	server.Deps.EventBus.Subscribe(event.ChannelCreatedEventType, func(evt event.Event) {
		fmt.Println("Channel created event received")
		// Cast the event to ChannelCreatedEvent
		channelCreatedEvent, ok := evt.(*event.ChannelCreatedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChannelCreatedEvent")
			return
		}

		fmt.Println("Channel created event received")

		logg := deps.Logger.With().
			Str("channelId", channelCreatedEvent.Channel.Id.String()).Logger()

		// Broadcast the created channels to all connected clients
		server.IterateClients(func(client *Client) bool {
			if client.CurrentSelectedWorkspace.Load() != channelCreatedEvent.Channel.WorkspaceId.String() {
				// Skip clients that are not in the same room
				return true
			}
			err := client.SendMessage(&outbound.OutboundChannelCreated{
				Channel: outbound.OutboundChannelCreatedChannel{
					Id:          channelCreatedEvent.Channel.Id,
					Name:        channelCreatedEvent.Channel.Name,
					Kind:        channelCreatedEvent.Channel.Kind,
					Topic:       channelCreatedEvent.Channel.Topic,
					IsPrivate:   channelCreatedEvent.Channel.IsPrivate,
					WorkspaceId: channelCreatedEvent.Channel.WorkspaceId,
					CreatedAt:   channelCreatedEvent.Channel.CreatedAt,
					UpdatedAt:   channelCreatedEvent.Channel.UpdatedAt,
					Index:       channelCreatedEvent.Channel.Index,
				},
			})

			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send channel create message to client")
				return false
			}
			return true
		})
	})

	server.Deps.EventBus.Subscribe(event.ChannelsReorderedEventType, func(evt event.Event) {
		// Cast the event to ChannelsReorderedEvent
		channelsReorderedEvent, ok := evt.(*event.ChannelsReorderedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChannelsReorderedEvent")
			return
		}

		// Convert []event.ChannelReorderMessage to []outbound.ChannelReorderMessage
		var outboundChannelReorders []outbound.ChannelReorderMessage
		for _, reorder := range channelsReorderedEvent.ChannelReorders {
			outboundChannelReorders = append(outboundChannelReorders, outbound.ChannelReorderMessage{
				ChannelId: reorder.ChannelId,
				NewOrder:  reorder.NewOrder,
			})
		}

		// Broadcast the reordered channels to all connected clients
		server.IterateClients(func(client *Client) bool {
			if client.CurrentSelectedWorkspace.Load() != channelsReorderedEvent.WorkspaceId.String() {
				// Skip clients that are not in the same room
				return true
			}
			err := client.SendMessage(&outbound.OutboundChannelsReordered{
				ChannelReorders: outboundChannelReorders,
			})
			if err != nil {
				server.Deps.Logger.Error().Err(err).
					Msg("failed to send channel reorder message to client")
				return false
			}
			return true
		})
	})

	server.Deps.EventBus.Subscribe(event.ChannelsDeletedEventType, func(evt event.Event) {
		// Cast the event to ChannelsDeletedEvent
		channelsDeletedEvent, ok := evt.(*event.ChannelsDeletedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChannelsDeletedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("channelId", channelsDeletedEvent.ChannelId.String()).Logger()
		// Broadcast the deleted channels to all connected clients
		server.IterateClients(func(client *Client) bool {
			if client.CurrentSelectedWorkspace.Load() != channelsDeletedEvent.WorkspaceId.String() {
				// Skip clients that are not in the same room
				return true
			}
			err := client.SendMessage(&outbound.OutboundChannelsDeleted{
				ChannelId:   channelsDeletedEvent.ChannelId,
				WorkspaceId: channelsDeletedEvent.WorkspaceId,
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send channel delete message to client")
				return false
			}
			return true
		})
	})

	server.Deps.EventBus.Subscribe(event.WorkspaceUpdatedEventType, func(evt event.Event) {
		workspaceUpdatedEvent, ok := evt.(*event.WorkspaceUpdatedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to WorkspaceUpdatedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("workspaceId", workspaceUpdatedEvent.Workspace.Id.String()).Logger()

		workspaceId := workspaceUpdatedEvent.Workspace.Id

		server.IterateClients(func(client *Client) bool {
			isMember, err := deps.IsUserInWorkspaceUseCase.Execute(
				context.Background(),
				workspaceId,
				client.UserId,
			)

			if err != nil {
				logg.Error().Err(err).
					Str("userId", client.UserId.String()).
					Msg("failed to check if user is in workspace")
				return true
			}

			if !isMember {
				return true
			}

			err = client.SendMessage(&outbound.OutboundWorkspaceUpdated{
				WorkspaceId: workspaceId.String(),
				Name:        workspaceUpdatedEvent.Workspace.Name,
				Topic:       workspaceUpdatedEvent.Workspace.Topic,
				Type:        string(workspaceUpdatedEvent.Workspace.Type),
			})
			if err != nil {
				logg.Error().Err(err).
					Str("userId", client.UserId.String()).
					Msg("failed to send workspace update message to client")
				return false
			}
			return true
		})
	})

	return server, nil
}

func (s *WsServer) Run() {
	pubsub := s.Deps.RedisClient.Client.Subscribe(context.Background(), "ws-messages")
	defer pubsub.Close()

	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)
		case client := <-s.Unregister:
			s.unregisterClient(client)
		case msg := <-pubsub.Channel():
			s.ForwardToClients([]byte(msg.Payload))
		}
	}
}

func (s *WsServer) ForwardToClients(message []byte) {
	for client := range s.clients {
		var forwardMessage ForwardMessage
		err := json.Unmarshal(message, &forwardMessage)
		if err != nil {
			s.Deps.Logger.Error().Err(err).Msg("Error on unmarshalling message")
			continue
		}

		if forwardMessage.EmitterServerId == s.backIdentifier {
			continue
		}

		s.Deps.Logger.Info().Str("message", string(message)).Msg("Forwarding message to client")

		client.send <- message
	}
}

func (s *WsServer) registerClient(client *Client) {
	s.clients[client] = true
}

func (s *WsServer) unregisterClient(client *Client) {
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
	}
}

func (s *WsServer) findRoomById(id string) *Room {
	var foundRoom *Room
	for iteratedRoom := range s.rooms {
		if iteratedRoom.Id == id {
			foundRoom = iteratedRoom
			break
		}
	}

	return foundRoom
}

func (s *WsServer) findClientByUserId(userId user_entity.UserId) *Client {
	for client := range s.IterateClients {
		if client.UserId == userId {
			return client
		}
	}

	return nil

}

func (s *WsServer) createRoom(name string, roomData RoomData) *Room {
	createdRoom := NewRoom(s.Deps, name, roomData)
	go createdRoom.RunRoom()
	s.rooms[createdRoom] = true

	return createdRoom
}

func (s *WsServer) IterateClients(fn func(client *Client) bool) {
	for client := range s.clients {
		if !fn(client) {
			break
		}
	}
}
