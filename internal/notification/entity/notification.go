package entity

import (
	group_message_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"time"
)

type NotificationId string
type NotificationType string

const (
	NotificationTypeDirectMessage   NotificationType = "direct_message"
	NotificationTypeChannelMessage  NotificationType = "channel_message"
	NotificationTypeWorkspaceInvite NotificationType = "workspace_invite"
	NotificationTypeGroupeMessage   NotificationType = "group_message"
)

type Notification struct {
	Id        NotificationId   `json:"id"`
	UserId    entity.UserId    `json:"user_id"`
	Type      NotificationType `json:"type"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
	// Metadata spécifiques selon le type
	DirectMessageData   *DirectMessageNotificationData   `json:"direct_message_data,omitempty"`
	ChannelMessageData  *ChannelMessageNotificationData  `json:"channel_message_data,omitempty"`
	WorkspaceInviteData *WorkspaceInviteNotificationData `json:"workspace_invite_data,omitempty"`
	GroupMessageData    *GroupMessageNotificationData    `json:"group_message_data,omitempty"`
}

type DirectMessageNotificationData struct {
	SenderId  entity.UserId
	MessageId chat_direct_entity.ChatDirectId
	//MessagePreview string
}

type ChannelMessageNotificationData struct {
	SenderId entity.UserId
	//SenderAvatarUrl string
	ChannelId   channel_entity.ChannelId
	WorkspaceId workspace_entity.WorkspaceId
	MessageId   channel_message_entity.ChannelMessageId
	//MessagePreview string
}

type WorkspaceInviteNotificationData struct {
	InviterId   entity.UserId
	WorkspaceId workspace_entity.WorkspaceId
}

type GroupMessageNotificationData struct {
	SenderId  entity.UserId
	GroupId   group_entity.GroupId
	MessageId group_message_entity.GroupChatMessageId
}

func (t NotificationType) String() string {
	return string(t)
}
