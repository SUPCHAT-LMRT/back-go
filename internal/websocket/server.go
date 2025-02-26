package websocket

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/websocket/room"
)

type WsServer struct {
	Deps       WebSocketDeps
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	rooms      map[*Room]bool
}

func NewWsServer(deps WebSocketDeps) *WsServer {
	return &WsServer{
		Deps:       deps,
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rooms:      make(map[*Room]bool),
	}
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

		// TODO impl forwardMessage.EmitterServerId
		if forwardMessage.EmitterServerId == "1" {
			continue
		}

		client.HandleNewMessage(forwardMessage.Payload)
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

func (s *WsServer) findClientById(clientId uuid.UUID) *Client {
	var foundClient *Client
	for client := range s.IterateClients {
		if client.Id == clientId {
			foundClient = client
			break
		}
	}

	return foundClient

}

func (s *WsServer) createRoom(name string, kind room.RoomKind) *Room {
	createdRoom := NewRoom(s.Deps, name, kind)
	go createdRoom.RunRoom()
	s.rooms[createdRoom] = true

	return createdRoom
}

func (s *WsServer) IterateClients(fn func(client *Client) bool) {
	for client := range s.clients {
		if fn(client) {
			break
		}
	}
}
