package reoder_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type ReorderChannelRequest struct {
	Id       string `json:"id"`
	NewIndex int    `json:"newIndex"`
}

type ReorderChannelHandler struct {
	reorderUseCase *ReorderChannelsUseCase
}

func NewReorderChannelHandler(reorderUseCase *ReorderChannelsUseCase) *ReorderChannelHandler {
	return &ReorderChannelHandler{reorderUseCase: reorderUseCase}
}

// Handle réorganise l'ordre des canaux dans un espace de travail
// @Summary Réorganisation des canaux
// @Description Modifie l'ordre d'affichage des canaux dans un espace de travail
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channels body []ReorderChannelRequest true "Liste des canaux à réordonner avec leur nouvel index"
// @Success 200 {string} string "Canaux réorganisés avec succès"
// @Failure 400 {object} map[string]string "Requête invalide ou mal formée"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour réorganiser les canaux"
// @Failure 404 {object} map[string]string "Canal non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la réorganisation des canaux"
// @Router /api/workspaces/{workspace_id}/channels/reorder [post]
// @Security ApiKeyAuth
func (h *ReorderChannelHandler) Handle(c *gin.Context) {
	var requests []ReorderChannelRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var inputs []ReorderChannelsInput
	for _, req := range requests {
		inputs = append(inputs, ReorderChannelsInput{
			ChannelId: entity.ChannelId(req.Id),
			NewIndex:  req.NewIndex,
		})
	}

	if err := h.reorderUseCase.ExecuteBulk(c, inputs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
