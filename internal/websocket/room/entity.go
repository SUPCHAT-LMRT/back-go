package room

type RoomKind string

// RoomKind specify the type of a room
const (
	ChannelRoomKind RoomKind = "channel"
	GroupRoomKind   RoomKind = "group"
	DirectRoomKind  RoomKind = "direct"
)
