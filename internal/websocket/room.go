package websocket

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type Room struct {
	deps       WebSocketDeps
	Id         string   `json:"id"`
	Data       RoomData `json:"data"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type RoomData interface {
	mustImplementRoomData()
}

type ChannelRoomData struct{}

func (ChannelRoomData) mustImplementRoomData() {}

type DirectRoomData struct {
	UserId      user_entity.UserId `json:"userId"`
	OtherUserId user_entity.UserId `json:"otherUserId"`
}

func (DirectRoomData) mustImplementRoomData() {}

func (d DirectRoomData) String() string {
	// create unique room name combined to the two IDs, the room name will be the same for both users
	// so the ids are ordered
	if d.UserId.IsAfter(d.OtherUserId) {
		return fmt.Sprintf("direct-%s_%s", d.UserId.String(), d.OtherUserId.String())
	}

	return fmt.Sprintf("direct-%s_%s", d.OtherUserId.String(), d.UserId.String())
}

// NewRoom creates a new Room
func NewRoom(deps WebSocketDeps, id string, data RoomData) *Room {
	return &Room{
		deps:       deps,
		Id:         id,
		Data:       data,
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
	delete(room.clients, client)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}

	backIdentifier, err := room.deps.GetBackIdentifierUseCase.Execute(context.Background())
	if err != nil {
		room.deps.Logger.Error().Err(err).Msg("Error on getting back identifier")
		return
	}

	forwardedMessage, err := json.Marshal(
		ForwardMessage{EmitterServerId: backIdentifier, Payload: message},
	)
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
