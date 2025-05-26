package get_list_invite_link

import (
	"net/http"

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

	c.JSON(http.StatusOK, inviteLinks)
}
