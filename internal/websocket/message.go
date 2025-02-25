package websocket

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"log"
	"time"
)

const OutboundSendMessageAction = "send-message"
const OutboundRoomJoinedAction = "room-joined"
const OutboundChannelCreatedAction = "channel-created"

const InboundJoinDirectRoomAction = "join-direct-room"
const InboundJoinGroupRoomAction = "join-group-room"
const InboundJoinChannelRoomAction = "join-channel-room"
const InboundLeaveRoomAction = "leave-room"
const InboundUnselectWorkspaceAction = "unselect-workspace"
const InboundSelectWorkspaceAction = "select-workspace"

type Message struct {
	Id      uuid.UUID `json:"id"`
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  *Room     `json:"target"`
	Sender  *Client   `json:"sender"`
	// Payload is a placeholder for any additional data that needs to be sent with the message, depending on the action.
	Payload   any       `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

type MessageSender interface{}

type WorkspaceMessageSender struct {
	UserId            user_entity.UserId       `json:"userId"`
	WorkspaceMemberId entity.WorkspaceMemberId `json:"workspaceMemberId"`
	WorkspacePseudo   string                   `json:"workspacePseudo"`
}

type GroupDirectMessageSender struct {
	UserId user_entity.UserId `json:"userId"`
}

func (m *Message) encode() []byte {
	result, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}

	return result
}
