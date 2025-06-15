package entity

import (
	"time"

	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type GroupChatMessageId string
type GroupChatAttachmentId string

type GroupChatMessage struct {
	Id          GroupChatMessageId
	GroupId     group_entity.GroupId
	Content     string
	AuthorId    entity.UserId
	Reactions   []*MessageReaction
	Attachments []*GroupChatMessageAttachment
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupChatMessageAttachment struct {
	Id       GroupChatAttachmentId
	FileName string
}

func (id GroupChatMessageId) String() string {
	return string(id)
}

func (id GroupChatAttachmentId) String() string {
	return string(id)
}
