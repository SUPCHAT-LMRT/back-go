package websocket

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/outbound"
	"github.com/supchat-lmrt/back-go/internal/websocket/room"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"strings"
	"sync/atomic"
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
	Id            uuid.UUID
	UserId        user_entity.UserId `json:"userId"`
	CurrentRoomId atomic.Value
	conn          *websocket.Conn
	wsServer      *WsServer
	rooms         map[*Room]bool
	send          chan []byte
}

func NewClient(user *user_entity.User, conn *websocket.Conn, wsServer *WsServer) *Client {
	c := &Client{
		Id:            uuid.New(),
		UserId:        user.Id,
		CurrentRoomId: atomic.Value{},
		conn:          conn,
		wsServer:      wsServer,
		rooms:         make(map[*Room]bool),
		send:          make(chan []byte, 256),
	}

	c.CurrentRoomId.Store("")

	return c
}

func (c *Client) HandleNewMessage(jsonMessage []byte) {
	var message messages.DefaultMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
		return
	}

	message.SetId(uuid.NewString())
	message.SetCreatedAt(time.Now())
	//message.SetEmittedBy(c)

	c.wsServer.Deps.Logger.Info().
		Str("action", string(message.Action)).
		Str("id", message.Id).
		Str("emittedBy", c.UserId.String()).
		Msg("New message")

	switch message.Action {
	case messages.InboundJoinChannelRoomAction:
		joinChannelMessage := inbound.InboundJoinChannel{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &joinChannelMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleJoinChannelRoomMessage(&joinChannelMessage)
		break
	case messages.InboundJoinDirectRoomAction:
		joinDirectRoomMessage := inbound.InboundJoinDirectRoom{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &joinDirectRoomMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleJoinDirectRoomMessage(&joinDirectRoomMessage)
		break
	case messages.InboundSendChannelMessageAction:
		sendMessage := inbound.InboundSendMessageToChannel{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &sendMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleSendMessageToChannel(&sendMessage)
		break
	case messages.InboundSendDirectMessageAction:
		sendMessage := inbound.InboundSendDirectMessage{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &sendMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleSendDirectMessage(&sendMessage)
		break
	case messages.InboundSelectWorkspaceAction:
		selectWorkspaceMessage := inbound.InboundSelectWorkspace{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &selectWorkspaceMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleSelectWorkspaceMessage(&selectWorkspaceMessage)
		break
	case messages.InboundUnselectWorkspaceAction:
		unselectWorkspaceMessage := inbound.InboundUnselectWorkspace{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &unselectWorkspaceMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleUnselectWorkspaceMessage(&unselectWorkspaceMessage)
		break
	case messages.InboundLeaveRoomAction:
		leaveRoomMessage := inbound.InboundLeaveRoom{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &leaveRoomMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleLeaveRoomMessage(&leaveRoomMessage)
		break
	case messages.InboundChannelMessageReactionToggle:
		reactionMessage := inbound.InboundChannelMessageReactionToggle{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &reactionMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleChannelMessageReactionToggleMessage(&reactionMessage)
		break
	case messages.InboundDirectMessageReactionToggle:
		reactionMessage := inbound.InboundDirectMessageReactionToggle{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &reactionMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleDirectMessageReactionToggleMessage(&reactionMessage)
		break
	default:
		log.Printf("Unknown action %s", message.Action)
	}
}

func (c *Client) handleSendMessageToChannel(message *inbound.InboundSendMessageToChannel) {
	if strings.TrimSpace(message.Content) == "" {
		return
	}

	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := message.ChannelId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	channelSender, err := c.toOutboundSendChannelMessageSender(roomId)
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		messageId := bson.NewObjectID().Hex()
		err = foundRoom.SendMessage(&outbound.OutboundSendMessageToChannel{
			MessageId: messageId,
			Content:   message.Content,
			ChannelId: message.ChannelId,
			Sender:    channelSender,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.SendChannelMessageObservers {
			observer.OnSendMessage(message, entity.ChannelMessageId(messageId), c.UserId)
		}
	}

}

func (c *Client) handleSendDirectMessage(message *inbound.InboundSendDirectMessage) {
	if strings.TrimSpace(message.Content) == "" {
		return
	}

	roomId := c.buildDirectMessageRoomId(message.OtherUserId)
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	messageSender, err := c.toOutboundDirectMessageSender()
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		messageId := bson.NewObjectID().Hex()
		err = foundRoom.SendMessage(&outbound.OutboundSendDirectMessage{
			MessageId:   messageId,
			Content:     message.Content,
			Sender:      messageSender,
			OtherUserId: message.OtherUserId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.SendDirectMessageObservers {
			observer.OnSendMessage(message, chat_direct_entity.ChatDirectId(messageId), c.UserId)
		}
	}
}

func (c *Client) handleJoinChannelRoomMessage(message *inbound.InboundJoinChannel) {
	// Todo: Check if the room exists
	c.joinRoom(message.ChannelId.String(), room.ChannelRoomKind)
}

func (c *Client) handleJoinDirectRoomMessage(message *inbound.InboundJoinDirectRoom) {
	roomId := c.buildDirectMessageRoomId(message.OtherUserId)
	// Todo: Check if the room exists
	c.joinRoom(roomId, room.DirectRoomKind)
}

//	func (c *Client) handleJoinGroupRoomMessage(message Message) {
//		roomId := message.Message
//		// Todo: Check if the room exists
//		c.joinRoom(roomId, GroupRoomKind, message.Sender)
//	}

func (c *Client) handleLeaveRoomMessage(message *inbound.InboundLeaveRoom) {
	roomId := message.RoomId

	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	if _, ok := c.rooms[foundRoom]; ok {
		delete(c.rooms, foundRoom)
	}

	foundRoom.unregister <- c
}

func (c *Client) handleSelectWorkspaceMessage(message *inbound.InboundSelectWorkspace) {
	c.CurrentRoomId.Store(message.WorkspaceId.String())
}

func (c *Client) handleUnselectWorkspaceMessage(message *inbound.InboundUnselectWorkspace) {
	c.CurrentRoomId.Store("")
}

//func (c *Client) handleJoinRoomPrivateMessage(message Message) {
//	clientId, err := uuid.Parse(message.Message)
//	if err != nil {
//		log.Println("Error parsing room id")
//		return
//	}
//
//	target := c.wsServer.findClientById(clientId)
//	if target == nil {
//		return
//	}
//
//	// create unique room name combined to the two IDs
//	roomName := message.Message + c.ChannelId.String()
//
//	c.joinRoom(roomName, DirectRoomKind, target)
//	target.joinRoom(roomName, DirectRoomKind, c)
//}

func (c *Client) handleChannelMessageReactionToggleMessage(message *inbound.InboundChannelMessageReactionToggle) {
	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := message.RoomId
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		member, err := c.toOutboundChannelMessageReactionMember(roomId)
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
			return
		}

		added, err := c.wsServer.Deps.ToggleReactionChannelMessageUseCase.Execute(context.Background(), entity.ChannelMessageId(message.MessageId), c.UserId, message.Reaction)
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating reaction")
			return
		}

		if added {
			err = foundRoom.SendMessage(&outbound.OutboundChannelMessageReactionAdded{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		} else {
			err = foundRoom.SendMessage(&outbound.OutboundChannelMessageReactionAdded{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		}
	}
}

func (c *Client) handleDirectMessageReactionToggleMessage(message *inbound.InboundDirectMessageReactionToggle) {
	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := c.buildDirectMessageRoomId(message.OtherUserId)
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		member, err := c.toOutboundDirectMessageReactionMember()
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
			return
		}

		added, err := c.wsServer.Deps.ToggleReactionDirectMessageUseCase.Execute(context.Background(), chat_direct_entity.ChatDirectId(message.MessageId), c.UserId, message.Reaction)
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating reaction")
			return
		}

		if added {
			err = foundRoom.SendMessage(&outbound.OutboundDirectMessageReactionAdded{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		} else {
			err = foundRoom.SendMessage(&outbound.OutboundDirectMessageReactionRemoved{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		}
	}
}

func (c *Client) joinRoom(roomId string, kind room.RoomKind) {
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		// Todo handle GroupRoomKind
		foundRoom = c.wsServer.createRoom(roomId, kind)
	}

	if !c.isInRoom(foundRoom) {
		c.rooms[foundRoom] = true
		foundRoom.register <- c
		c.notifyRoomJoined(foundRoom)
	}
}

func (c *Client) isInRoom(room *Room) bool {
	if _, ok := c.rooms[room]; ok {
		return true
	}
	return false
}

func (c *Client) SendMessage(message messages.Message) error {
	encoded, err := message.Encode()
	if err != nil {
		return err
	}

	c.send <- encoded
	// TODO Notify on redis

	return nil
}

func (c *Client) notifyRoomJoined(room *Room) {
	message := outbound.OutboundRoomJoined{
		Room: outbound.OutboundRoomJoinedRoom{
			Id:   room.Id,
			Kind: room.Kind,
		},
	}

	encoded, err := message.Encode()
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on encoding message")
		return
	}

	c.send <- encoded
}

func (c *Client) buildDirectMessageRoomId(otherUserId user_entity.UserId) string {
	// create unique room name combined to the two IDs, the room name will be the same for both users
	// so the ids are ordered
	if c.UserId.IsAfter(otherUserId) {
		return fmt.Sprintln("direct-", c.UserId.String(), "_", otherUserId.String())
	} else {
		return fmt.Sprintln("direct-", otherUserId.String(), "_", c.UserId.String())
	}
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

		c.HandleNewMessage(jsonMessage)
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
	for iteratedRoom := range c.rooms {
		iteratedRoom.unregister <- c
	}
	close(c.send)
	c.conn.Close()
}

func (c *Client) toOutboundSendChannelMessageSender(roomId string) (*outbound.OutboundSendMessageToChannelSender, error) {
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

	username := workspaceMember.Pseudo
	if username == "" {
		username = user.FullName()
	}

	return &outbound.OutboundSendMessageToChannelSender{
		UserId:            user.Id,
		Pseudo:            user.FullName(),
		WorkspaceMemberId: workspaceMember.Id,
		WorkspacePseudo:   username,
	}, nil
}

func (c *Client) toOutboundChannelMessageReactionMember(roomId string) (*outbound.OutboundChannelMessageReactionMember, error) {
	channel, err := c.wsServer.Deps.GetChannelUseCase.Execute(context.Background(), channel_entity.ChannelId(roomId))
	if err != nil {
		return nil, err
	}

	workspaceMember, err := c.wsServer.Deps.GetWorkspaceMemberUseCase.Execute(context.Background(), channel.WorkspaceId, c.UserId)
	if err != nil {
		return nil, err
	}

	username := workspaceMember.Pseudo
	// If the user is not in the workspace, we will use the user pseudo (fallback)
	if username == "" {
		user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
		if err != nil {
			return nil, err
		}
		username = user.FullName()
	}

	return &outbound.OutboundChannelMessageReactionMember{
		UserId:   c.UserId.String(),
		Username: username,
	}, nil
}

func (c *Client) toOutboundDirectMessageReactionMember() (*outbound.OutboundDirectMessageReactionMember, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundDirectMessageReactionMember{
		UserId:   c.UserId.String(),
		Username: user.FullName(),
	}, nil
}

func (c *Client) toOutboundDirectMessageSender() (*outbound.OutboundSendDirectMessageSender, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundSendDirectMessageSender{
		UserId:    user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
