package websocket

import (
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/websocket/room"
)

type WsServer struct {
	Deps       WebSocketDeps
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewWsServer(deps WebSocketDeps) *WsServer {
	return &WsServer{
		Deps:       deps,
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

func (s *WsServer) Run() {
	//pubsub := s.Deps.RedisClient.Client.Subscribe(context.Background(), "ws-messages")
	//defer pubsub.Close()

	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)
		case client := <-s.Unregister:
			s.unregisterClient(client)
		case message := <-s.Broadcast:
			s.broadcastToClients(message)
			//case msg := <-pubsub.Channel():
			//	s.ForwardToClients([]byte(msg.Payload))
		}
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

func (s *WsServer) broadcastToClients(message []byte) {
	for client := range s.clients {
		client.send <- message
	}
}

func (s *WsServer) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range s.rooms {
		if room.Id == id {
			foundRoom = room
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
