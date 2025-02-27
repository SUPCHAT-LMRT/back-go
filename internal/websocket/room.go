package websocket

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/outbound"
	"github.com/supchat-lmrt/back-go/internal/websocket/room"
)

type Room struct {
	deps       WebSocketDeps
	Id         string        `json:"id"`
	Kind       room.RoomKind `json:"kind"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

// NewRoom creates a new Room
func NewRoom(deps WebSocketDeps, id string, kind room.RoomKind) *Room {
	return &Room{
		deps:       deps,
		Id:         id,
		Kind:       kind,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}

	// TODO impl forwardMessage.EmitterServerId
	forwardedMessage, err := json.Marshal(ForwardMessage{EmitterServerId: "1", Payload: message})
	if err != nil {
		room.deps.Logger.Error().Err(err).Msg("Error on forwarding message to clients")
		return
	}

	room.deps.RedisClient.Client.Publish(context.Background(), "ws-messages", forwardedMessage)
}

func (room *Room) SendMessage(message messages.Message) error {
	encoded, err := message.Encode()
	if err != nil {
		return err
	}

	room.broadcastToClientsInRoom(encoded)
	return nil
}

func (room *Room) SendChannelMessage(message outbound.OutboundSendMessageToChannel) error {
	encoded, err := message.Encode()
	if err != nil {
		return err
	}

	room.broadcastToClientsInRoom(encoded)
	return nil
}
