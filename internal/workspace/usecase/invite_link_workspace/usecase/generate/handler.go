package generate

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"net/http"
)

type CreateInviteLinkHandler struct {
	usecase *InviteLinkUseCase
}

func NewCreateInviteLinkHandler(usecase *InviteLinkUseCase) *CreateInviteLinkHandler {
	return &CreateInviteLinkHandler{usecase: usecase}
}

func (h *CreateInviteLinkHandler) Handle(c *gin.Context) {

	inviteLink, err := h.usecase.CreateInviteLink(c, entity.WorkspaceId(c.Param("workspaceId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, inviteLink)
}
