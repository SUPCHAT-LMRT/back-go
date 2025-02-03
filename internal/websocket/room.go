package websocket

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type RoomKind string

const (
	ChannelRoomKind RoomKind = "channel"
	GroupRoomKind   RoomKind = "group"
	DirectRoomKind  RoomKind = "direct"
)

type Room struct {
	deps       WebSocketDeps
	Id         string   `json:"id"`
	Kind       RoomKind `json:"kind"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

// NewRoom creates a new Room
func NewRoom(deps WebSocketDeps, id string, kind RoomKind) *Room {
	return &Room{
		deps:       deps,
		Id:         id,
		Kind:       kind,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
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

		case message := <-room.broadcast:
			// If the user is not in the room he broadcasted to, don't send the message.
			if message.Sender.isInRoom(room) {
				room.broadcastToClientsInRoom(message.encode())
			}
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	// by sending the message first the new user won't see his own message.
	if room.Kind != DirectRoomKind {
		room.notifyClientJoined(client)
	}
	room.clients[client] = true

	room.restoreMessages(client)
}

func (room *Room) restoreMessages(client *Client) {
	switch room.Kind {
	case ChannelRoomKind:
		room.restoreChannelMessages(client)
	case GroupRoomKind:
		room.restoreGroupMessages(client)
	case DirectRoomKind:
		room.restoreDirectMessages(client)
	}
}

func (room *Room) restoreChannelMessages(client *Client) {
	channelMessages, err := room.deps.ListChannelMessagesUseCase.Execute(context.Background(), channel_entity.ChannelId(room.Id))
	if err != nil {
		return
	}

	for _, message := range channelMessages {
		user, err := room.deps.GetUserByIdUseCase.Execute(context.Background(), message.AuthorId)
		if err != nil {
			continue
		}

		channel, err := room.deps.GetChannelUseCase.Execute(context.Background(), channel_entity.ChannelId(room.Id))
		if err != nil {
			continue
		}

		member, err := room.deps.GetWorkspaceMemberUseCase.Execute(context.Background(), channel.WorkspaceId, message.AuthorId)
		if err != nil {
			continue
		}

		client.SendMessage(Message{
			Id:      uuid.New(),
			Action:  SendMessageAction,
			Message: message.Content,
			Target:  room,
			MessageSender: &WorkspaceMessageSender{
				UserId:            message.AuthorId,
				Pseudo:            user.Pseudo,
				WorkspaceMemberId: member.Id,
				WorkspacePseudo:   member.Pseudo,
			},
		})
	}
}

func (room *Room) restoreGroupMessages(client *Client) {
	fmt.Println("restoreGroupMessages")
}

func (room *Room) restoreDirectMessages(client *Client) {
	fmt.Println("restoreDirectMessages")
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
	
	// by sending the message after the user won't see his own message.
	if room.Kind != DirectRoomKind {
		room.notifyClientLeft(client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *Client) {
	room.SendMessage(Message{
		Id:      uuid.New(),
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.Id),
	})
}

const goodbyeMessage = "%s left the room"

func (room *Room) notifyClientLeft(client *Client) {
	room.SendMessage(Message{
		Id:      uuid.New(),
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(goodbyeMessage, client.Id),
	})
}

func (room *Room) SendMessage(message Message) {
	room.broadcastToClientsInRoom(message.encode())
}
