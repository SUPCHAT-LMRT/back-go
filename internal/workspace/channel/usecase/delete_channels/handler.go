package delete_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type DeleteChannelHandler struct {
	useCase *DeleteChannelUseCase
}

func NewDeleteChannelHandler(useCase *DeleteChannelUseCase) *DeleteChannelHandler {
	return &DeleteChannelHandler{useCase: useCase}
}

// Handle supprime un canal dans un espace de travail
// @Summary Suppression d'un canal
// @Description Supprime un canal existant dans l'espace de travail spécifié
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channel_id path string true "ID du canal à supprimer"
// @Success 204 {string} string "Canal supprimé avec succès"
// @Failure 400 {object} map[string]string "ID de canal manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour supprimer le canal"
// @Failure 404 {object} map[string]string "Canal non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la suppression du canal"
// @Router /api/workspaces/{workspace_id}/channels/{channel_id} [delete]
// @Security ApiKeyAuth
func (h *DeleteChannelHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	if channelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel_id is required"})
		return
	}

	err := h.useCase.Execute(c, entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
