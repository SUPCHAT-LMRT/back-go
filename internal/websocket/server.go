package websocket

import (
	"context"
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
