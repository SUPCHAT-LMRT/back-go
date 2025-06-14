package messages

import (
	"time"

	"github.com/google/uuid"
)

type Action string

// Outbound actions are actions that are sent from the server to the client.
const (
	OutboundSendChannelMessageAction            Action = "send-channel-message"
	OutboundSendDirectMessageAction             Action = "send-direct-message"
	OutboundSendGroupMessageAction              Action = "send-group-message"
	OutboundChannelRoomJoinedAction             Action = "channel-room-joined"
	OutboundDirectRoomJoinedAction              Action = "direct-room-joined"
	OutboundGroupRoomJoinedAction               Action = "group-room-joined"
	OutboundChannelCreatedAction                Action = "channel-created"
	OutboundChannelMessageReactionAddedAction   Action = "channel-message-reaction-added"
	OutboundChannelMessageReactionRemovedAction Action = "channel-message-reaction-removed"
	OutboundDirectMessageReactionAddedAction    Action = "direct-message-reaction-added"
	OutboundDirectMessageReactionRemovedAction  Action = "direct-message-reaction-removed"
	OutboundGroupMessageReactionAddedAction     Action = "group-message-reaction-added"
	OutboundGroupMessageReactionRemovedAction   Action = "group-message-reaction-removed"
	OutboundGroupMessageContentEditedAction     Action = "group-message-content-edited"
	OutboundGroupOwnershipTransferredAction     Action = "group-ownership-transferred"
	OutboundGroupMessageDeletedAction           Action = "group-message-deleted"
	OutboundRecentDirectChatAddedAction         Action = "recent-direct-chat-added"
	OutboundRecentGroupChatAddedAction          Action = "recent-group-chat-added"
	OutboundRecentGroupChatRemovedAction        Action = "recent-group-chat-removed"
	OutboundGroupMemberAddedAction              Action = "group-member-added"
	OutboundGroupMemberRemovedAction            Action = "group-member-removed"
	OutboundUserStatusUpdatedAction             Action = "user-status-updated"
	OutboundSelfStatusUpdatedAction             Action = "self-status-updated"
	OutboundChannelsReorderedAction             Action = "channels-reordered"
	OutboundChannelsDeletedAction               Action = "channels-deleted"
	OutboundWorkspaceUpdatedAction              Action = "workspace-updated"
)

// Inbound actions are actions that are sent from the client to the server.
const (
	InboundSendChannelMessageAction     Action = "send-channel-message"
	InboundSendDirectMessageAction      Action = "send-direct-message"
	InboundSendGroupMessageAction       Action = "send-group-message"
	InboundJoinDirectRoomAction         Action = "join-direct-room"
	InboundJoinGroupRoomAction          Action = "join-group-room"
	InboundJoinChannelRoomAction        Action = "join-channel-room"
	InboundLeaveRoomAction              Action = "leave-room"
	InboundUnselectWorkspaceAction      Action = "unselect-workspace"
	InboundSelectWorkspaceAction        Action = "select-workspace"
	InboundChannelMessageReactionToggle Action = "channel-message-reaction-toggle"
	InboundDirectMessageReactionToggle  Action = "direct-message-reaction-toggle"
	InboundGroupMessageReactionToggle   Action = "group-message-reaction-toggle"
	InboundGroupMessageContentEdit      Action = "group-message-content-edit"
	InboundGroupMessageDelete           Action = "group-message-delete"
)

type Message interface {
	GetActionName() Action
	SetId(string)
	SetCreatedAt(time.Time)
	Encode() ([]byte, error)
	mustExtendDefaultMessage()
}

type DefaultMessage struct {
	TransportMessageId        string    `json:"transportMessageId"`
	Action                    Action    `json:"action"`
	TransportMessageCreatedAt time.Time `json:"transportMessageCreatedAt"`
}

func (m *DefaultMessage) SetId(id string) {
	m.TransportMessageId = id
}

func (m *DefaultMessage) SetCreatedAt(createdAt time.Time) {
	m.TransportMessageCreatedAt = createdAt
}

//nolint:unused
func (m *DefaultMessage) mustExtendDefaultMessage() {}

func NewDefaultMessage(action Action) DefaultMessage {
	return DefaultMessage{
		TransportMessageId:        uuid.NewString(),
		Action:                    action,
		TransportMessageCreatedAt: time.Now(),
	}
}
