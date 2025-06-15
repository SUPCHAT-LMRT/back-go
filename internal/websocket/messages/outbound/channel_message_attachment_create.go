package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"time"
)

type OutboundChannelMessageAttachmentCreated struct {
	messages.DefaultMessage
	Message *OutboundChannelMessageAttachmentCreatedMessage `json:"message"`
}

type OutboundChannelMessageAttachmentCreatedMessage struct {
	Id                      string    `json:"id"`
	SenderUserId            string    `json:"senderUserId"`
	SenderPseudo            string    `json:"senderPseudo"`
	SenderWorkspaceMemberId string    `json:"senderWorkspaceMemberId"`
	AttachmentFileId        string    `json:"attachmentFileId"`
	AttachmentFileName      string    `json:"attachmentFileName"`
	CreatedAt               time.Time `json:"createdAt"`
}

func (o *OutboundChannelMessageAttachmentCreated) GetActionName() messages.Action {
	return messages.OutboundChannelMessageAttachmentCreatedAction
}

func (o *OutboundChannelMessageAttachmentCreated) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
