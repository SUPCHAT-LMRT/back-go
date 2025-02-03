package websocket

type SendMessageObserver interface {
	OnSendMessage(message Message)
}
