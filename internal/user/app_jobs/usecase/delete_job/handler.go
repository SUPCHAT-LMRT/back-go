package delete_job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
)

type DeleteJobHandler struct {
	useCase *DeleteJobUseCase
}

func NewDeleteJobHandler(useCase *DeleteJobUseCase) *DeleteJobHandler {
	return &DeleteJobHandler{useCase: useCase}
}

// Handle supprime un rôle de travail du système
// @Summary Supprimer un rôle de travail
// @Description Supprime un rôle de travail existant en fonction de son ID
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "ID du rôle à supprimer"
// @Success 200 {object} map[string]string "Rôle supprimé avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/{id} [delete]
// @Security ApiKeyAuth
func (h *DeleteJobHandler) Handle(c *gin.Context) {
	jobId := c.Param("id")
	if jobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	err := h.useCase.Execute(c.Request.Context(), entity.JobId(jobId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}
