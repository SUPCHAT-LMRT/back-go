package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"time"
)

type OutboundChatDirectAttachmentCreated struct {
	messages.DefaultMessage
	Message *OutboundChatDirectAttachmentCreatedMessage `json:"message"`
}

type OutboundChatDirectAttachmentCreatedMessage struct {
	Id                 string    `json:"id"`
	AuthorUserId       string    `json:"authorUserId"`
	AuthorFirstName    string    `json:"authorFirstName"`
	AuthorLastName     string    `json:"authorLastName"`
	AttachmentFileId   string    `json:"attachmentFileId"`
	AttachmentFileName string    `json:"attachmentFileName"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (o *OutboundChatDirectAttachmentCreated) GetActionName() messages.Action {
	return messages.OutboundDirectMessageAttachmentCreatedAction
}

func (o *OutboundChatDirectAttachmentCreated) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
