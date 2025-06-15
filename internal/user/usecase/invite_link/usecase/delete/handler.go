package delete

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteInviteLinkHandler struct {
	usecase *DeleteInviteLinkUseCase
}

func NewDeleteInviteLinkHandler(usecase *DeleteInviteLinkUseCase) *DeleteInviteLinkHandler {
	return &DeleteInviteLinkHandler{usecase: usecase}
}

// Handle supprime un lien d'invitation
// @Summary Supprimer un lien d'invitation
// @Description Supprime définitivement un lien d'invitation du système
// @Tags account
// @Accept json
// @Produce json
// @Param token path string true "Token unique du lien d'invitation à supprimer"
// @Success 204 "Lien d'invitation supprimé avec succès"
// @Failure 400 {object} map[string]string "Token manquant"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/invite-link/{token} [delete]
// @Security ApiKeyAuth
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
