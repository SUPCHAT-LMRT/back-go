package websocket

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"log"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type Client struct {
	Id       uuid.UUID
	UserId   user_entity.UserId `json:"userId"`
	conn     *websocket.Conn
	wsServer *WsServer
	rooms    map[*Room]bool
	send     chan []byte
}

func NewClient(user *user_entity.User, conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		Id:       uuid.New(),
		UserId:   user.Id,
		conn:     conn,
		wsServer: wsServer,
		rooms:    make(map[*Room]bool),
		send:     make(chan []byte, 256),
	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	message.Id = uuid.New()
	// Attach the client object as the sender of the messsage.
	message.Sender = c

	switch message.Action {
	case SendMessageAction:
		c.handleSendMessage(message)
	case JoinDirectRoomAction:
		c.handleJoinDirectRoomMessage(message)
	case JoinGroupRoomAction:
		c.handleJoinGroupRoomMessage(message)
	case JoinChannelRoomAction:
		c.handleJoinChannelRoomMessage(message)
	case LeaveRoomAction:
		c.handleLeaveRoomMessage(message)
	case JoinRoomPrivateAction:
		c.handleJoinRoomPrivateMessage(message)
	}
}

func (c *Client) handleSendMessage(message Message) {
	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := message.Target.Id
	room := c.wsServer.findRoomById(roomId)
	if room == nil {
		return
	}

	var err error
	message.MessageSender, err = c.toMessageSender(room)
	if err != nil {
		log.Printf("Error getting sender: %s", err)
		return
	}
	// Use the ChatServer method to find the room, and if found, broadcast!
	if room = c.wsServer.findRoomById(roomId); room != nil {
		room.SendMessage(message)
	}

	// Notify observers
	for _, observer := range c.wsServer.Deps.SendMessageObservers {
		observer.OnSendMessage(message)
	}
}

func (c *Client) handleJoinDirectRoomMessage(message Message) {
	roomId := message.Message
	// Todo: Check if the room exists
	c.joinRoom(roomId, DirectRoomKind, message.Sender)
}

func (c *Client) handleJoinGroupRoomMessage(message Message) {
	roomId := message.Message
	// Todo: Check if the room exists
	c.joinRoom(roomId, GroupRoomKind, message.Sender)
}

func (c *Client) handleJoinChannelRoomMessage(message Message) {
	roomId := message.Message
	// Todo: Check if the room exists
	c.joinRoom(roomId, ChannelRoomKind, message.Sender)
}

func (c *Client) handleLeaveRoomMessage(message Message) {
	roomId := message.Message

	room := c.wsServer.findRoomById(roomId)
	if room == nil {
		return
	}

	if _, ok := c.rooms[room]; ok {
		delete(c.rooms, room)
	}

	room.unregister <- c
}

func (c *Client) handleJoinRoomPrivateMessage(message Message) {
	clientId, err := uuid.Parse(message.Message)
	if err != nil {
		log.Println("Error parsing room id")
		return
	}

	target := c.wsServer.findClientById(clientId)
	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	roomName := message.Message + c.Id.String()

	c.joinRoom(roomName, DirectRoomKind, target)
	target.joinRoom(roomName, DirectRoomKind, c)
}

func (c *Client) joinRoom(roomId string, kind RoomKind, sender *Client) {
	room := c.wsServer.findRoomById(roomId)
	if room == nil {
		// Todo handle GroupRoomKind
		room = c.wsServer.createRoom(roomId, kind)
	}

	if !c.isInRoom(room) {
		c.rooms[room] = true
		room.register <- c
		c.notifyRoomJoined(room, sender)
	}
}

func (c *Client) isInRoom(room *Room) bool {
	if _, ok := c.rooms[room]; ok {
		return true
	}
	return false
}

func (c *Client) SendMessage(message Message) {
	c.send <- message.encode()
}

func (c *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Id:     uuid.New(),
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	c.send <- message.encode()
}

func (c *Client) ReadPump() {
	defer c.disconnect()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.handleNewMessage(jsonMessage)
	}
}

var (
	newline = []byte{'\n'}
)

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				fmt.Println("Error on close writer", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) disconnect() {
	c.wsServer.Unregister <- c
	for room := range c.rooms {
		room.unregister <- c
	}
	close(c.send)
	c.conn.Close()
}

func (c *Client) toMessageSender(room *Room) (MessageSender, error) {
	if room.Kind == ChannelRoomKind {
		return c.toWorkspaceMessageSender(room.Id)
	}

	return c.toGroupDirectMessageSender()
}

func (c *Client) toWorkspaceMessageSender(roomId string) (MessageSender, error) {
	// The room id is the channel id
	channel, err := c.wsServer.Deps.GetChannelUseCase.Execute(context.Background(), channel_entity.ChannelId(roomId))
	if err != nil {
		return nil, err
	}

	workspaceMember, err := c.wsServer.Deps.GetWorkspaceMemberUseCase.Execute(context.Background(), channel.WorkspaceId, c.UserId)
	if err != nil {
		return nil, err
	}

	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceMessageSender{
		UserId:            user.Id,
		Pseudo:            user.Pseudo,
		WorkspaceMemberId: workspaceMember.Id,
		WorkspacePseudo:   workspaceMember.Pseudo,
	}, nil
}

func (c *Client) toGroupDirectMessageSender() (MessageSender, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &GroupDirectMessageSender{
		UserId: user.Id,
		Pseudo: user.Pseudo,
	}, nil
}
