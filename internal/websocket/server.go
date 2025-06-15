package websocket

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/group/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
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

//nolint:revive
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
			user2, err := deps.GetUserByIdUseCase.Execute(
				context.Background(),
				messageSavedEvent.Message.User2Id,
			)
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
			user1, err := deps.GetUserByIdUseCase.Execute(
				context.Background(),
				messageSavedEvent.Message.User1Id,
			)
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

	server.Deps.EventBus.Subscribe(event.GroupCreatedEventType, func(evt event.Event) {
		groupCreatedEvent, ok := evt.(*event.GroupCreatedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to DirectChatMessageSavedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("groupId", groupCreatedEvent.Message.Id.String()).
			Str("groupName", groupCreatedEvent.Message.Name).Logger()

		// Broadcast the created group to all the invited users
		groupMembers, err := server.Deps.ListGroupMembersUseCase.Execute(context.Background(), groupCreatedEvent.Message.Id)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to list group members")
			return
		}

		for _, member := range groupMembers {
			client := server.findClientByUserId(member.UserId)
			if client == nil {
				logg.Warn().Str("userId", member.UserId.String()).
					Msg("client not found for group member, skipping")
				continue
			}

			err = client.SendMessage(&outbound.OutboundAddRecentGroupChat{
				GroupId:  groupCreatedEvent.Message.Id,
				ChatName: groupCreatedEvent.Message.Name,
			})
			if err != nil {
				logg.Error().Err(err).
					Str("userId", member.UserId.String()).
					Msg("failed to send message to group member")
				return
			}
		}
	})

	server.Deps.EventBus.Subscribe(event.GroupMemberAddedEventType, func(evt event.Event) {
		ctx := context.Background()
		groupMemberAddedEvent, ok := evt.(*event.GroupMemberAddedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to GroupMemberAddedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("groupId", groupMemberAddedEvent.Message.Id.String()).
			Str("userId", groupMemberAddedEvent.InvitedUserId.String()).Logger()

		// Send to the invited user, the group chat
		client := server.findClientByUserId(groupMemberAddedEvent.InvitedUserId)
		if client == nil {
			logg.Warn().Msg("client not found for invited user, skipping")
			return
		}

		err = client.SendMessage(&outbound.OutboundAddRecentGroupChat{
			GroupId:  groupMemberAddedEvent.Message.Id,
			ChatName: groupMemberAddedEvent.Message.Name,
		})
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to send message to invited user")
			return
		}

		// Broadcast the created group to all the invited users that the group member was added
		groupMembers, err := server.Deps.ListGroupMembersUseCase.Execute(ctx, groupMemberAddedEvent.Message.Id)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to list group members")
			return
		}

		invitedUserMember, err := server.Deps.GetMemberByUserUseCase.Execute(ctx, groupMemberAddedEvent.Message.Id, groupMemberAddedEvent.InvitedUserId)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to get invited user member")
			return
		}

		invitedUser, err := server.Deps.GetUserByIdUseCase.Execute(ctx, invitedUserMember.UserId)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to get invited user")
			return
		}

		invitedUserStatus, err := server.Deps.GetPublicStatusUseCase.Execute(ctx, invitedUserMember.UserId, entity.StatusOffline)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to get invited user status")
			return
		}

		for _, member := range groupMembers {
			// Skip the invited user, as they already received the message and fetch the group chat
			if member.UserId == groupMemberAddedEvent.InvitedUserId {
				continue
			}

			memberClient := server.findClientByUserId(member.UserId)
			if memberClient == nil {
				logg.Warn().Str("userId", member.UserId.String()).
					Msg("client not found for group member, skipping")
				continue
			}

			err = memberClient.SendMessage(&outbound.OutboundAddGroupMemberChat{
				GroupId: groupMemberAddedEvent.Message.Id,
				Member: &outbound.OutboundAddGroupMemberChatMember{
					Id:           invitedUserMember.Id,
					UserId:       invitedUserMember.UserId,
					UserName:     invitedUser.FullName(),
					IsGroupOwner: false,
					Status:       invitedUserStatus,
				},
			})
		}
	})

	server.Deps.EventBus.Subscribe(event.GroupMemberRemovedEventType, func(evt event.Event) {
		groupMemberRemovedEvent, ok := evt.(*event.GroupMemberRemovedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to GroupMemberRemovedEvent")
			return
		}

		logg := deps.Logger.With().
			Str("groupId", groupMemberRemovedEvent.Group.Id.String()).
			Str("userId", groupMemberRemovedEvent.RemovedMemberId.String()).Logger()

		// Broadcast the removed group member to all the group members
		groupMembers, err := server.Deps.ListGroupMembersUseCase.Execute(context.Background(), groupMemberRemovedEvent.Group.Id)
		if err != nil {
			if errors.Is(err, repository.ErrGroupNotFound) {
				// Send message to the removed user to remove the recent chat
				client := server.findClientByUserId(groupMemberRemovedEvent.RemovedUserId)
				if client == nil {
					logg.Warn().Str("userId", groupMemberRemovedEvent.RemovedUserId.String()).
						Msg("client not found for group member, skipping")
					return
				}

				err = client.SendMessage(&outbound.OutboundRemovedRecentGroupChat{
					GroupId: groupMemberRemovedEvent.Group.Id,
				})
				if err != nil {
					logg.Error().Err(err).
						Msg("failed to send message to invited user")
					return
				}
				return
			}
			logg.Error().Err(err).
				Msg("failed to list group members")
			return
		}

		// Send message to the removed user to remove the recent chat
		client := server.findClientByUserId(groupMemberRemovedEvent.RemovedUserId)
		if client == nil {
			logg.Warn().Str("userId", groupMemberRemovedEvent.RemovedUserId.String()).
				Msg("client not found for group member, skipping")
			return
		}

		err = client.SendMessage(&outbound.OutboundRemovedRecentGroupChat{
			GroupId: groupMemberRemovedEvent.Group.Id,
		})
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to send message to invited user")
			return
		}

		// Send message to all other group members about the removed member
		for _, member := range groupMembers {
			client := server.findClientByUserId(member.UserId)
			if client == nil {
				logg.Warn().Str("userId", member.UserId.String()).
					Msg("client not found for group member, skipping")
				continue
			}

			err = client.SendMessage(&outbound.OutboundRemoveGroupMemberChat{
				GroupId:  groupMemberRemovedEvent.Group.Id,
				MemberId: groupMemberRemovedEvent.RemovedMemberId,
				UserId:   groupMemberRemovedEvent.RemovedUserId,
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send message to group member")
				return
			}
		}
	})

	server.Deps.EventBus.Subscribe(event.GroupTransferOwnershipEventType, func(evt event.Event) {
		groupTransferOwnershipEvent, ok := evt.(*event.GroupTransferOwnershipEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to GroupTransferOwnershipEvent")
			return
		}

		logg := deps.Logger.With().
			Str("groupId", groupTransferOwnershipEvent.Group.Id.String()).
			Str("newOwnerId", groupTransferOwnershipEvent.NewOwnerId.String()).Logger()

		// Broadcast the ownership transfer to all group members
		groupMembers, err := server.Deps.ListGroupMembersUseCase.Execute(context.Background(), groupTransferOwnershipEvent.Group.Id)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to list group members")
			return
		}

		for _, member := range groupMembers {
			client := server.findClientByUserId(member.UserId)
			if client == nil {
				logg.Warn().Str("userId", member.UserId.String()).
					Msg("client not found for group member, skipping")
				continue
			}

			err = client.SendMessage(&outbound.OutboundGroupOwnershipTransferer{
				GroupId:    groupTransferOwnershipEvent.Group.Id,
				NewOwnerId: groupTransferOwnershipEvent.NewOwnerId,
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send ownership transfer message to group member")
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
		// Cast the event to ChannelCreatedEvent
		channelCreatedEvent, ok := evt.(*event.ChannelCreatedEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChannelCreatedEvent")
			return
		}

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
			outboundChannelReorders = append(
				outboundChannelReorders,
				outbound.ChannelReorderMessage{
					ChannelId: reorder.ChannelId,
					NewOrder:  reorder.NewOrder,
				},
			)
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

	server.Deps.EventBus.Subscribe(event.ChannelAttachmentSentEventType, func(evt event.Event) {
		channelAttachmentSentEvent, ok := evt.(*event.ChannelAttachmentSentEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChannelAttachmentSentEvent")
			return
		}

		if len(channelAttachmentSentEvent.ChannelMessage.Attachments) == 0 {
			server.Deps.Logger.Error().Msg("channelAttachmentSentEvent has no attachments")
			return
		}

		logg := deps.Logger.With().
			Str("channelId", channelAttachmentSentEvent.ChannelMessage.ChannelId.String()).
			Str("attachmentId", channelAttachmentSentEvent.ChannelMessage.Attachments[0].Id.String()).Logger()

		server.IterateClients(func(client *Client) bool {
			if client.CurrentSelectedWorkspace.Load() != channelAttachmentSentEvent.WorkspaceId.String() {
				// Skip clients that are not in the same room
				return true
			}

			user, err := deps.GetUserByIdUseCase.Execute(context.Background(), channelAttachmentSentEvent.ChannelMessage.AuthorId)
			if err != nil {
				logg.Error().Err(err).
					Str("userId", channelAttachmentSentEvent.ChannelMessage.AuthorId.String()).
					Msg("failed to get user by ID")
				return false
			}

			err = client.SendMessage(&outbound.OutboundChannelMessageAttachmentCreated{
				Message: &outbound.OutboundChannelMessageAttachmentCreatedMessage{
					Id:                      channelAttachmentSentEvent.ChannelMessage.Id.String(),
					SenderUserId:            channelAttachmentSentEvent.ChannelMessage.AuthorId.String(),
					SenderPseudo:            user.FullName(),
					SenderWorkspaceMemberId: channelAttachmentSentEvent.WorkspaceMemberId.String(),
					AttachmentFileId:        channelAttachmentSentEvent.ChannelMessage.Attachments[0].Id.String(),
					AttachmentFileName:      channelAttachmentSentEvent.ChannelMessage.Attachments[0].FileName,
					CreatedAt:               channelAttachmentSentEvent.ChannelMessage.CreatedAt,
				},
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send attachment sent message to client")
				return false
			}
			return true
		})
	})

	server.Deps.EventBus.Subscribe(event.ChatDirectAttachmentSentEventType, func(evt event.Event) {
		chatDirectAttachmentSentEvent, ok := evt.(*event.ChatDirectAttachmentSentEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to ChatDirectAttachmentSentEvent")
			return
		}

		if len(chatDirectAttachmentSentEvent.ChatDirect.Attachments) == 0 {
			server.Deps.Logger.Error().Msg("chatDirectAttachmentSentEvent has no attachments")
			return
		}

		logg := deps.Logger.With().
			Str("otherUserId", chatDirectAttachmentSentEvent.ChatDirect.GetReceiverId().String()).
			Str("attachmentId", chatDirectAttachmentSentEvent.ChatDirect.Attachments[0].Id.String()).Logger()

		// Send to both users in the chat direct attachment sent event (sender and receiver)
		for _, iteratedUserId := range []user_entity.UserId{chatDirectAttachmentSentEvent.ChatDirect.SenderId, chatDirectAttachmentSentEvent.ChatDirect.GetReceiverId()} {
			client := server.findClientByUserId(iteratedUserId)
			if client == nil {
				logg.Warn().Str("userId", iteratedUserId.String()).
					Msg("client not found for user, skipping")
				return
			}

			authorUser, err := deps.GetUserByIdUseCase.Execute(context.Background(), chatDirectAttachmentSentEvent.ChatDirect.SenderId)
			if err != nil {
				logg.Error().Err(err).
					Str("userId", chatDirectAttachmentSentEvent.ChatDirect.SenderId.String()).
					Msg("failed to get user by ID")
				return
			}

			err = client.SendMessage(&outbound.OutboundChatDirectAttachmentCreated{
				Message: &outbound.OutboundChatDirectAttachmentCreatedMessage{
					Id:                 chatDirectAttachmentSentEvent.ChatDirect.Id.String(),
					AuthorUserId:       authorUser.Id.String(),
					AuthorFirstName:    authorUser.FirstName,
					AuthorLastName:     authorUser.LastName,
					AttachmentFileId:   chatDirectAttachmentSentEvent.ChatDirect.Attachments[0].Id.String(),
					AttachmentFileName: chatDirectAttachmentSentEvent.ChatDirect.Attachments[0].FileName,
					CreatedAt:          chatDirectAttachmentSentEvent.ChatDirect.CreatedAt,
				},
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send attachment sent message to client")
				return
			}
		}
	})

	server.Deps.EventBus.Subscribe(event.GroupAttachmentSentEventType, func(evt event.Event) {
		groupAttachmentSentEvent, ok := evt.(*event.GroupAttachmentSentEvent)
		if !ok {
			server.Deps.Logger.Error().Msg("failed to cast event to GroupAttachmentSentEvent")
			return
		}

		if len(groupAttachmentSentEvent.GroupChatMessage.Attachments) == 0 {
			server.Deps.Logger.Error().Msg("groupAttachmentSentEvent has no attachments")
			return
		}

		logg := deps.Logger.With().
			Str("groupId", groupAttachmentSentEvent.GroupChatMessage.GroupId.String()).
			Str("attachmentId", groupAttachmentSentEvent.GroupChatMessage.Attachments[0].Id.String()).Logger()

		groupMembers, err := server.Deps.ListGroupMembersUseCase.Execute(context.Background(), groupAttachmentSentEvent.GroupChatMessage.GroupId)
		if err != nil {
			logg.Error().Err(err).
				Msg("failed to list group members")
			return
		}

		for _, groupMember := range groupMembers {
			client := server.findClientByUserId(groupMember.UserId)
			if client == nil {
				logg.Warn().Str("groupMemberId", groupMember.Id.String()).
					Str("userId", groupMember.UserId.String()).
					Msg("client not found for user, skipping")
				return
			}

			authorUser, err := deps.GetUserByIdUseCase.Execute(context.Background(), groupAttachmentSentEvent.GroupChatMessage.AuthorId)
			if err != nil {
				logg.Error().Err(err).
					Str("userId", groupAttachmentSentEvent.GroupChatMessage.AuthorId.String()).
					Msg("failed to get user by ID")
				return
			}

			err = client.SendMessage(&outbound.OutboundGroupAttachmentCreated{
				GroupId: groupAttachmentSentEvent.GroupChatMessage.GroupId,
				Message: &outbound.OutboundGroupAttachmentCreatedMessage{
					Id:                 groupAttachmentSentEvent.GroupChatMessage.Id.String(),
					AuthorUserId:       authorUser.Id.String(),
					AuthorFirstName:    authorUser.FirstName,
					AuthorLastName:     authorUser.LastName,
					AttachmentFileId:   groupAttachmentSentEvent.GroupChatMessage.Attachments[0].Id.String(),
					AttachmentFileName: groupAttachmentSentEvent.GroupChatMessage.Attachments[0].FileName,
					CreatedAt:          groupAttachmentSentEvent.GroupChatMessage.CreatedAt,
				},
			})
			if err != nil {
				logg.Error().Err(err).
					Msg("failed to send attachment sent message to client")
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
	delete(s.clients, client)
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
