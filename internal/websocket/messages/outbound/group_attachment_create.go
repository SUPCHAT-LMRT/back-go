package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"time"
)

type OutboundGroupAttachmentCreated struct {
	messages.DefaultMessage
	GroupId group_entity.GroupId                   `json:"groupId"`
	Message *OutboundGroupAttachmentCreatedMessage `json:"message"`
}

type OutboundGroupAttachmentCreatedMessage struct {
	Id                 string    `json:"id"`
	AuthorUserId       string    `json:"authorUserId"`
	AuthorFirstName    string    `json:"authorFirstName"`
	AuthorLastName     string    `json:"authorLastName"`
	AttachmentFileId   string    `json:"attachmentFileId"`
	AttachmentFileName string    `json:"attachmentFileName"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (o *OutboundGroupAttachmentCreated) GetActionName() messages.Action {
	return messages.OutboundGroupMessageAttachmentCreatedAction
}

func (o *OutboundGroupAttachmentCreated) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
