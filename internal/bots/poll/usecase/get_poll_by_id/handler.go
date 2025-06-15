package get_poll_by_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPollByIdHandler struct {
	usecase *GetPollByIdUseCase
}

// Handle récupère un sondage par son identifiant
// @Summary Récupérer un sondage
// @Description Récupère les détails d'un sondage existant dans l'espace de travail
// @Tags poll
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param poll_id path string true "ID du sondage à récupérer"
// @Success 200 {object} interface{} "Détails du sondage"
// @Failure 400 {object} map[string]string "ID du sondage manquant"
// @Failure 404 {object} map[string]string "Sondage non trouvé"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Router /api/workspaces/{workspace_id}/poll/{poll_id} [get]
// @Security ApiKeyAuth
func NewGetPollByIdHandler(usecase *GetPollByIdUseCase) *GetPollByIdHandler {
	return &GetPollByIdHandler{usecase: usecase}
}

func (h *GetPollByIdHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poll_id is required"})
		return
	}

	poll, err := h.usecase.Execute(c, pollId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, poll)
}
