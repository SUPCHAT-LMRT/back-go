package room

type RoomKind string

const (
	ChannelRoomKind RoomKind = "channel"
	GroupRoomKind   RoomKind = "group"
	DirectRoomKind  RoomKind = "direct"
)
