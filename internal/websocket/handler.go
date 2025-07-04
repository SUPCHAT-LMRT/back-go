package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type WebsocketHandler struct {
	wsServer *WsServer
}

func NewWebsocketHandler(wsServer *WsServer) *WebsocketHandler {
	return &WebsocketHandler{wsServer: wsServer}
}

func (h *WebsocketHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(user, conn, h.wsServer)

	go client.WritePump()
	go client.ReadPump()

	h.wsServer.Register <- client
}
