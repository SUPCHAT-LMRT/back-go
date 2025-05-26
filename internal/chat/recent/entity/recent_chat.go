package entity

import "time"

type (
	RecentChatKind int
	RecentChatId   string
)

const (
	RecentChatKindGroup RecentChatKind = iota
	RecentChatKindDirect
)

type RecentChat struct {
	Id        RecentChatId
	Kind      RecentChatKind
	Name      string
	UpdatedAt time.Time
}
