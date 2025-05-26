package save_status

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
)

type SaveStatusHandler struct {
	useCase *SaveStatusUseCase
}

func NewSaveStatusHandler(useCase *SaveStatusUseCase) *SaveStatusHandler {
	return &SaveStatusHandler{useCase: useCase}
}

func (h *SaveStatusHandler) Handle(c *gin.Context) {
	user := c.MustGet("user").(*entity.User) //nolint:revive

	var body struct {
		Status string `json:"status" binding:"required,userStatus"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.Execute(c, user.Id, user_status_entity.Status(body.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
