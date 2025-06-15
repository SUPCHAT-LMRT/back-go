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

// Handle récupère la liste des liens d'invitation créés
// @Summary Lister les liens d'invitation
// @Description Récupère tous les liens d'invitation actifs avec leurs informations
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {array} get_list_invite_link.ResponseInviteLink "Liste des liens d'invitation"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 403 {object} map[string]string "Accès refusé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/invite-link [get]
// @Security ApiKeyAuth
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
