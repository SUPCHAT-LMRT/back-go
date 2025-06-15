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

// Handle récupère les données associées à un lien d'invitation
// @Summary Obtenir les données d'un lien d'invitation
// @Description Récupère les informations associées à un lien d'invitation spécifique
// @Tags account
// @Accept json
// @Produce json
// @Param token path string true "Token unique du lien d'invitation"
// @Success 200 {object} get_data_token_invite.InviteLinkDataResponse "Données de l'invitation"
// @Failure 400 {object} map[string]string "Token manquant"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/invite-link/{token} [get]
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
