package get_list_invite_link

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetListInviteLinkHandler struct {
	usecase *GetListInviteLinkUseCase
}

func NewGetListInviteLinkHandler(usecase *GetListInviteLinkUseCase) *GetListInviteLinkHandler {
	return &GetListInviteLinkHandler{usecase: usecase}
}

func (h *GetListInviteLinkHandler) Handle(c *gin.Context) {
	inviteLinks, err := h.usecase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseInviteLinks []ResponseInviteLink
	for _, inviteLink := range inviteLinks {
		responseInviteLinks = append(responseInviteLinks, ResponseInviteLink{
			Token:          inviteLink.Token,
			Email:          inviteLink.Email,
			ExpirationDate: inviteLink.ExpiresAt,
		})
	}

	c.JSON(http.StatusOK, responseInviteLinks)
}

type ResponseInviteLink struct {
	Token          string    `json:"token"`
	Email          string    `json:"email"`
	ExpirationDate time.Time `json:"expiresAt"`
}
