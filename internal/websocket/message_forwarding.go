package websocket

type ForwardMessage struct {
	EmitterServerId string
	Payload         []byte
}
