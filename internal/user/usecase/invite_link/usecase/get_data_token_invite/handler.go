package get_data_token_invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetInviteLinkDataHandler struct {
	usecase *GetInviteLinkDataUseCase
}

func NewGetInviteLinkDataHandler(usecase *GetInviteLinkDataUseCase) *GetInviteLinkDataHandler {
	return &GetInviteLinkDataHandler{usecase: usecase}
}

func (h *GetInviteLinkDataHandler) Handle(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	inviteLink, err := h.usecase.GetInviteLinkData(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, InviteLinkDataResponse{
		FirstName: inviteLink.FirstName,
		LastName:  inviteLink.LastName,
		Email:     inviteLink.Email,
	})
}

type InviteLinkDataResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
