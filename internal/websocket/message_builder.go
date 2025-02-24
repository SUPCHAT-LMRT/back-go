package websocket

import (
	"github.com/google/uuid"
	"time"
)

type MessageBuilder struct {
	Action    string
	Message   string
	Target    *Room
	Sender    *Client
	Payload   any
	CreatedAt time.Time
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (m *MessageBuilder) WithAction(action string) *MessageBuilder {
	m.Action = action
	return m
}

func (m *MessageBuilder) WithMessage(message string) *MessageBuilder {
	m.Message = message
	return m
}

func (m *MessageBuilder) WithTarget(target *Room) *MessageBuilder {
	m.Target = target
	return m
}

func (m *MessageBuilder) WithSender(sender *Client) *MessageBuilder {
	m.Sender = sender
	return m
}

func (m *MessageBuilder) WithPayload(payload any) *MessageBuilder {
	m.Payload = payload
	return m
}

func (m *MessageBuilder) WithCreatedAt(date time.Time) *MessageBuilder {
	m.CreatedAt = date
	return m
}

func (m *MessageBuilder) Build() Message {
	return Message{
		Id:        uuid.New(),
		Action:    m.Action,
		Message:   m.Message,
		Target:    m.Target,
		Sender:    m.Sender,
		Payload:   m.Payload,
		CreatedAt: m.CreatedAt,
	}
}
