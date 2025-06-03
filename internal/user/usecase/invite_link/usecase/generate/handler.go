package generate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateInviteLinkRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
}

type CreateInviteLinkHandler struct {
	usecase *InviteLinkUseCase
}

func NewCreateInviteLinkHandler(usecase *InviteLinkUseCase) *CreateInviteLinkHandler {
	return &CreateInviteLinkHandler{usecase: usecase}
}

func (h *CreateInviteLinkHandler) Handle(c *gin.Context) {
	var req CreateInviteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviteLink, err := h.usecase.CreateInviteLink(c, req.FirstName, req.LastName, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, inviteLink)
}
