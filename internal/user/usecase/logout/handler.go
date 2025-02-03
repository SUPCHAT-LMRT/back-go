package logout

import (
	"github.com/gin-gonic/gin"
	"os"
)

type LogoutHandler struct {
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

func (l LogoutHandler) Handle(c *gin.Context) {
	c.SetCookie("accessToken", "", -1, "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refreshToken", "", -1, "/", os.Getenv("DOMAIN"), false, true)
}
