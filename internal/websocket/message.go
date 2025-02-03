package websocket

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"log"
)

const SendMessageAction = "send-message"
const JoinDirectRoomAction = "join-direct-room"
const JoinGroupRoomAction = "join-group-room"
const JoinChannelRoomAction = "join-channel-room"
const LeaveRoomAction = "leave-room"
const UserConnectAction = "user-connect"
const UserDisconnectAction = "user-disconnect"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const ChannelCreatedAction = "channel-created"

type Message struct {
	Id      uuid.UUID `json:"id"`
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  *Room     `json:"target"`
	Sender  *Client   `json:"sender"`
	// Payload is a placeholder for any additional data that needs to be sent with the message, depending on the action.
	Payload any `json:"payload"`
}

type MessageSender interface{}

type WorkspaceMessageSender struct {
	UserId            user_entity.UserId       `json:"userId"`
	Pseudo            string                   `json:"pseudo"`
	WorkspaceMemberId entity.WorkspaceMemberId `json:"workspaceMemberId"`
	WorkspacePseudo   string                   `json:"workspacePseudo"`
}

type GroupDirectMessageSender struct {
	UserId user_entity.UserId `json:"userId"`
	Pseudo string             `json:"pseudo"`
}

func (m *Message) encode() []byte {
	result, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}

	return result
}
