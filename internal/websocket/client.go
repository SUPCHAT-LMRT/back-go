package websocket

import (
	"context"
	"fmt"
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/toggle_reaction"
	"log"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/outbound"
	"github.com/supchat-lmrt/back-go/internal/websocket/utils"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	Id                       uuid.UUID
	UserId                   user_entity.UserId `json:"userId"`
	CurrentSelectedWorkspace atomic.Value
	conn                     *websocket.Conn
	wsServer                 *WsServer
	rooms                    map[*Room]bool
	send                     chan []byte
}

func NewClient(user *user_entity.User, conn *websocket.Conn, wsServer *WsServer) *Client {
	c := &Client{
		Id:                       uuid.New(),
		UserId:                   user.Id,
		CurrentSelectedWorkspace: atomic.Value{},
		conn:                     conn,
		wsServer:                 wsServer,
		rooms:                    make(map[*Room]bool),
		send:                     make(chan []byte, 256),
	}

	c.CurrentSelectedWorkspace.Store("")

	return c
}

//nolint:revive
func (c *Client) HandleNewMessage(jsonMessage []byte) {
	var message messages.DefaultMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
		return
	}

	message.SetId(uuid.NewString())
	message.SetCreatedAt(time.Now())
	// message.SetEmittedBy(c)

	c.wsServer.Deps.Logger.Debug().
		Str("action", string(message.Action)).
		Str("id", message.TransportMessageId).
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
	case messages.InboundJoinDirectRoomAction:
		joinDirectRoomMessage := inbound.InboundJoinDirectRoom{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &joinDirectRoomMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleJoinDirectRoomMessage(&joinDirectRoomMessage)
	case messages.InboundJoinGroupRoomAction:
		joinGroupRoomMessage := inbound.InboundJoinGroupRoom{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &joinGroupRoomMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleJoinGroupRoomMessage(&joinGroupRoomMessage)
	case messages.InboundSendChannelMessageAction:
		sendMessage := inbound.InboundSendMessageToChannel{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &sendMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleSendMessageToChannel(&sendMessage)
	case messages.InboundSendDirectMessageAction:
		sendMessage := inbound.InboundSendDirectMessage{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &sendMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleSendDirectMessage(&sendMessage)
	case messages.InboundSendGroupMessageAction:
		sendMessage := inbound.InboundSendGroupMessage{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &sendMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleSendGroupMessage(&sendMessage)
	case messages.InboundSelectWorkspaceAction:
		selectWorkspaceMessage := inbound.InboundSelectWorkspace{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &selectWorkspaceMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleSelectWorkspaceMessage(&selectWorkspaceMessage)
	case messages.InboundUnselectWorkspaceAction:
		unselectWorkspaceMessage := inbound.InboundUnselectWorkspace{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &unselectWorkspaceMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleUnselectWorkspaceMessage(&unselectWorkspaceMessage)
	case messages.InboundLeaveRoomAction:
		leaveRoomMessage := inbound.InboundLeaveRoom{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &leaveRoomMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleLeaveRoomMessage(&leaveRoomMessage)
	case messages.InboundChannelMessageReactionToggle:
		reactionMessage := inbound.InboundChannelMessageReactionToggle{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &reactionMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}

		c.handleChannelMessageReactionToggleMessage(&reactionMessage)
	case messages.InboundDirectMessageReactionToggle:
		reactionMessage := inbound.InboundDirectMessageReactionToggle{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &reactionMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleDirectMessageReactionToggleMessage(&reactionMessage)
	case messages.InboundGroupMessageReactionToggle:
		reactionMessage := inbound.InboundGroupMessageReactionToggle{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &reactionMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleGroupMessageReactionToggleMessage(&reactionMessage)
	case messages.InboundGroupMessageContentEdit:
		editMessage := inbound.InboundGroupMessageContentEdit{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &editMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleGroupMessageContentEdit(&editMessage)
	case messages.InboundGroupMessageDelete:
		deleteMessage := inbound.InboundGroupMessageDelete{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &deleteMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleGroupMessageDelete(&deleteMessage)
	case messages.InboundDirectMessageContentEdit:
		editMessage := inbound.InboundDirectMessageContentEdit{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &editMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleDirectMessageContentEdit(&editMessage)
	case messages.InboundDirectMessageDelete:
		deleteMessage := inbound.InboundDirectMessageDelete{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &deleteMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleDirectMessageDelete(&deleteMessage)
	case messages.InboundChannelMessageContentEdit:
		editMessage := inbound.InboundChannelMessageContentEdit{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &editMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleChannelMessageContentEdit(&editMessage)
	case messages.InboundChannelMessageDelete:
		deleteMessage := inbound.InboundChannelMessageDelete{DefaultMessage: message}
		if err := json.Unmarshal(jsonMessage, &deleteMessage); err != nil {
			log.Printf("Error on unmarshal JSON message %s %s", err, string(jsonMessage))
			return
		}
		c.handleChannelMessageDelete(&deleteMessage)
	default:
		log.Printf("Unknown action %s", message.Action)
	}
}

//nolint:revive
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
			CreatedAt: message.TransportMessageCreatedAt,
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

//nolint:revive
func (c *Client) handleSendDirectMessage(message *inbound.InboundSendDirectMessage) {
	if strings.TrimSpace(message.Content) == "" {
		return
	}

	roomId := utils.BuildDirectMessageRoomId(c.UserId, message.OtherUserId)
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
			CreatedAt:   message.TransportMessageCreatedAt,
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

//nolint:revive
func (c *Client) handleSendGroupMessage(message *inbound.InboundSendGroupMessage) {
	if strings.TrimSpace(message.Content) == "" {
		return
	}

	roomId := message.GroupId
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	messageSender, err := c.toOutboundGroupMessageSender()
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		messageId := bson.NewObjectID().Hex()
		err = foundRoom.SendMessage(&outbound.OutboundSendGroupMessage{
			MessageId: messageId,
			Content:   message.Content,
			GroupId:   message.GroupId,
			Sender:    messageSender,
			CreatedAt: message.TransportMessageCreatedAt,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.SendGroupMessageObservers {
			observer.OnSendMessage(message, group_chat_entity.GroupChatMessageId(messageId), c.UserId)
		}
	}
}

func (c *Client) handleJoinChannelRoomMessage(message *inbound.InboundJoinChannel) {
	c.joinRoom(message.ChannelId.String(), &ChannelRoomData{})
}

func (c *Client) handleJoinDirectRoomMessage(message *inbound.InboundJoinDirectRoom) {
	roomId := utils.BuildDirectMessageRoomId(c.UserId, message.OtherUserId)
	c.joinRoom(roomId, &DirectRoomData{UserId: c.UserId, OtherUserId: message.OtherUserId})
}

func (c *Client) handleJoinGroupRoomMessage(message *inbound.InboundJoinGroupRoom) {
	roomId := message.GroupId.String()
	c.joinRoom(roomId, &GroupRoomData{})
}

func (c *Client) handleLeaveRoomMessage(message *inbound.InboundLeaveRoom) {
	roomId := message.RoomId

	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	delete(c.rooms, foundRoom)

	foundRoom.unregister <- c
}

func (c *Client) handleSelectWorkspaceMessage(message *inbound.InboundSelectWorkspace) {
	c.CurrentSelectedWorkspace.Store(message.WorkspaceId.String())
}

func (c *Client) handleUnselectWorkspaceMessage(message *inbound.InboundUnselectWorkspace) {
	c.CurrentSelectedWorkspace.Store("")
}

// func (c *Client) handleJoinRoomPrivateMessage(message Message) {
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
// }

//nolint:revive
func (c *Client) handleChannelMessageReactionToggleMessage(
	message *inbound.InboundChannelMessageReactionToggle,
) {
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

		added, err := c.wsServer.Deps.ToggleReactionChannelMessageUseCase.Execute(
			context.Background(),
			entity.ChannelMessageId(message.MessageId),
			c.UserId,
			message.Reaction,
		)
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
			err = foundRoom.SendMessage(&outbound.OutboundChannelMessageReactionRemoved{
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

//nolint:revive
func (c *Client) handleDirectMessageReactionToggleMessage(
	message *inbound.InboundDirectMessageReactionToggle,
) {
	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := utils.BuildDirectMessageRoomId(c.UserId, message.OtherUserId)
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

		added, err := c.wsServer.Deps.ToggleReactionDirectMessageUseCase.Execute(
			context.Background(),
			chat_direct_entity.ChatDirectId(message.MessageId),
			c.UserId,
			message.Reaction,
		)
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating reaction")
			return
		}

		if added {
			err = foundRoom.SendMessage(&outbound.OutboundDirectMessageReactionAdded{
				MessageId:   message.MessageId,
				Reaction:    message.Reaction,
				OtherUserId: message.OtherUserId.String(),
				Member:      *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		} else {
			err = foundRoom.SendMessage(&outbound.OutboundDirectMessageReactionRemoved{
				MessageId:   message.MessageId,
				Reaction:    message.Reaction,
				OtherUserId: message.OtherUserId.String(),
				Member:      *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		}
	}
}

func (c *Client) handleGroupMessageReactionToggleMessage(
	message *inbound.InboundGroupMessageReactionToggle,
) {
	// The send-message action, this will send messages to a specific room now.
	// Which room wil depend on the message Target
	roomId := message.GroupId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		member, err := c.toOutboundGroupMessageReactionMember()
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating sender")
			return
		}

		added, err := c.wsServer.Deps.ToggleGroupChatReactionUseCase.Execute(
			context.Background(),
			toggle_reaction.ToggleReactionInput{
				MessageId: group_chat_entity.GroupChatMessageId(message.MessageId),
				UserId:    c.UserId,
				Reaction:  message.Reaction,
			},
		)
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on creating reaction")
			return
		}

		if added {
			err = foundRoom.SendMessage(&outbound.OutboundGroupMessageReactionAdded{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				GroupId:   message.GroupId,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		} else {
			err = foundRoom.SendMessage(&outbound.OutboundGroupMessageReactionRemoved{
				MessageId: message.MessageId,
				Reaction:  message.Reaction,
				GroupId:   message.GroupId,
				Member:    *member,
			})
			if err != nil {
				c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
				return
			}
		}
	}
}

func (c *Client) handleGroupMessageContentEdit(
	message *inbound.InboundGroupMessageContentEdit,
) {
	roomId := message.GroupId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundGroupMessageContentEdited{
			MessageId:  message.MessageId,
			NewContent: message.NewContent,
			GroupId:    message.GroupId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.EditGroupMessageObservers {
			observer.OnEditMessage(message)
		}
	}
}

func (c *Client) handleGroupMessageDelete(
	message *inbound.InboundGroupMessageDelete,
) {
	roomId := message.GroupId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundGroupMessageDeleted{
			MessageId: message.MessageId,
			GroupId:   message.GroupId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.DeleteGroupMessageObservers {
			observer.OnDeleteMessage(message)
		}
	}
}

func (c *Client) handleDirectMessageContentEdit(
	message *inbound.InboundDirectMessageContentEdit,
) {
	roomId := utils.BuildDirectMessageRoomId(c.UserId, message.OtherUserId)
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundDirectMessageContentEdited{
			OtherUserId: message.OtherUserId,
			MessageId:   message.MessageId,
			NewContent:  message.NewContent,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.EditDirectMessageObservers {
			observer.OnEditMessage(message)
		}
	}
}

func (c *Client) handleDirectMessageDelete(
	message *inbound.InboundDirectMessageDelete,
) {
	roomId := utils.BuildDirectMessageRoomId(c.UserId, message.OtherUserId)
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundDirectMessageDeleted{
			OtherUserId: message.OtherUserId,
			MessageId:   message.MessageId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.DeleteDirectMessageObservers {
			observer.OnDeleteMessage(message)
		}
	}
}

func (c *Client) handleChannelMessageContentEdit(
	message *inbound.InboundChannelMessageContentEdit,
) {
	roomId := message.ChannelId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundChannelMessageContentEdited{
			MessageId:  message.MessageId,
			NewContent: message.NewContent,
			ChannelId:  message.ChannelId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.EditChannelMessageObservers {
			observer.OnEditMessage(message)
		}
	}
}

func (c *Client) handleChannelMessageDelete(
	message *inbound.InboundChannelMessageDelete,
) {
	roomId := message.ChannelId.String()
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		return
	}

	// Use the ChatServer method to find the room, and if found, broadcast!
	if foundRoom = c.wsServer.findRoomById(roomId); foundRoom != nil {
		err := foundRoom.SendMessage(&outbound.OutboundChannelMessageDeleted{
			MessageId: message.MessageId,
			ChannelId: message.ChannelId,
		})
		if err != nil {
			c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on sending message")
			return
		}

		// Notify observers
		for _, observer := range c.wsServer.Deps.DeleteChannelMessageObservers {
			observer.OnDeleteMessage(message)
		}
	}
}

func (c *Client) joinRoom(roomId string, roomData RoomData) {
	foundRoom := c.wsServer.findRoomById(roomId)
	if foundRoom == nil {
		foundRoom = c.wsServer.createRoom(roomId, roomData)
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

func (c *Client) notifyRoomJoined(r *Room) {
	var message messages.Message
	switch roomData := r.Data.(type) {
	case *ChannelRoomData:
		message = &outbound.OutboundChannelRoomJoined{
			RoomId: r.Id,
		}
	case *DirectRoomData:
		message = &outbound.OutboundDirectRoomJoined{
			RoomId:      r.Id,
			OtherUserId: roomData.OtherUserId,
		}
	case *GroupRoomData:
		message = &outbound.OutboundGroupRoomJoined{
			RoomId: r.Id,
		}
	}

	if message == nil {
		c.wsServer.Deps.Logger.Warn().
			Any("data", r.Data).
			Str("data_type", reflect.TypeOf(r.Data).String()).
			Msg("Unknown room data")
		return
	}

	encoded, err := message.Encode()
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on encoding message")
		return
	}

	c.send <- encoded
}

func (c *Client) ReadPump() {
	defer c.disconnect()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Set the user status to offline
	err := c.wsServer.Deps.SaveStatusUseCase.Execute(
		context.Background(),
		c.UserId,
		user_status_entity.StatusOnline,
	)
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on setting user status to online")
		return
	}

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.HandleNewMessage(jsonMessage)
	}
}

var newline = []byte{'\n'}

//nolint:revive
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				fmt.Println("Error on close writer", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
	_ = c.conn.Close()

	// Set the user status to offline
	err := c.wsServer.Deps.SaveStatusUseCase.Execute(
		context.Background(),
		c.UserId,
		user_status_entity.StatusOffline,
	)
	if err != nil {
		c.wsServer.Deps.Logger.Error().Err(err).Msg("Error on setting user status to offline")
		return
	}
}

func (c *Client) toOutboundSendChannelMessageSender(
	roomId string,
) (*outbound.OutboundSendMessageToChannelSender, error) {
	// The room id is the channel id
	channel, err := c.wsServer.Deps.GetChannelUseCase.Execute(
		context.Background(),
		channel_entity.ChannelId(roomId),
	)
	if err != nil {
		return nil, err
	}

	workspaceMember, err := c.wsServer.Deps.GetWorkspaceMemberUseCase.Execute(
		context.Background(),
		channel.WorkspaceId,
		c.UserId,
	)
	if err != nil {
		return nil, err
	}

	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundSendMessageToChannelSender{
		UserId:            user.Id,
		Pseudo:            user.FullName(),
		WorkspaceMemberId: workspaceMember.Id,
	}, nil
}

func (c *Client) toOutboundChannelMessageReactionMember(
	roomId string,
) (*outbound.OutboundChannelMessageReactionMember, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(
		context.Background(),
		c.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundChannelMessageReactionMember{
		UserId:   c.UserId.String(),
		Username: user.FullName(),
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

func (c *Client) toOutboundGroupMessageReactionMember() (*outbound.OutboundGroupMessageReactionMember, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundGroupMessageReactionMember{
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

func (c *Client) toOutboundGroupMessageSender() (*outbound.OutboundSendGroupMessageSender, error) {
	user, err := c.wsServer.Deps.GetUserByIdUseCase.Execute(context.Background(), c.UserId)
	if err != nil {
		return nil, err
	}

	return &outbound.OutboundSendGroupMessageSender{
		UserId:    user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
