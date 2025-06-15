package delete_poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeletePollHandler struct {
	usecase *DeletePollUseCase
}

func NewDeletePollHandler(usecase *DeletePollUseCase) *DeletePollHandler {
	return &DeletePollHandler{usecase: usecase}
}

// Handle supprime un sondage existant
// @Summary Supprimer un sondage
// @Description Supprime un sondage existant dans l'espace de travail
// @Tags poll
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param poll_id path string true "ID du sondage à supprimer"
// @Success 200 {object} map[string]string "Sondage supprimé avec succès"
// @Failure 400 {object} map[string]string "ID du sondage manquant"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/workspaces/{workspace_id}/poll/{poll_id} [delete]
// @Security ApiKeyAuth
func (h *DeletePollHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poll_id is required"})
		return
	}

	err := h.usecase.Execute(c, pollId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poll deleted successfully"})
}
