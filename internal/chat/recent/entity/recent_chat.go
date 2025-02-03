package entity

import "time"

type RecentChatKind int
type RecentChatId string

const (
	RecentChatKindGroup RecentChatKind = iota
	RecentChatKindDirect
)

type RecentChat struct {
	Id        RecentChatId
	Kind      RecentChatKind
	AvatarUrl string
	Name      string
	UpdatedAt time.Time
}
