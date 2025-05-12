package delete

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteInviteLinkHandler struct {
	usecase *DeleteInviteLinkUseCase
}

func NewDeleteInviteLinkHandler(usecase *DeleteInviteLinkUseCase) *DeleteInviteLinkHandler {
	return &DeleteInviteLinkHandler{usecase: usecase}
}

func (h *DeleteInviteLinkHandler) Handle(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := h.usecase.Execute(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
