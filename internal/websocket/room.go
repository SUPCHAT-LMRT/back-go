package websocket

type RoomKind string

const (
	ChannelRoomKind RoomKind = "channel"
	GroupRoomKind   RoomKind = "group"
	DirectRoomKind  RoomKind = "direct"
)

type Room struct {
	deps       WebSocketDeps
	Id         string   `json:"id"`
	Kind       RoomKind `json:"kind"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

// NewRoom creates a new Room
func NewRoom(deps WebSocketDeps, id string, kind RoomKind) *Room {
	return &Room{
		deps:       deps,
		Id:         id,
		Kind:       kind,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			// If the user is not in the room he broadcasted to, don't send the message.
			if message.Sender.isInRoom(room) {
				room.broadcastToClientsInRoom(message.encode())
			}
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}
func (room *Room) SendMessage(message Message) {
	room.broadcastToClientsInRoom(message.encode())
}
