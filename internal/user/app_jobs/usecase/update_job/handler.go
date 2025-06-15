package update_job

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models" // Import pour que Swagger trouve les modèles
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"net/http"
)

type UpdateJobHandler struct {
	useCase *UpdateJobUseCase
}

func NewUpdateJobHandler(useCase *UpdateJobUseCase) *UpdateJobHandler {
	return &UpdateJobHandler{useCase: useCase}
}

// Handle met à jour un rôle de travail existant
// @Summary Mettre à jour un rôle de travail
// @Description Met à jour le nom d'un rôle de travail existant
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "ID du rôle à mettre à jour"
// @Param request body models.UpdateJobRequest true "Informations pour mettre à jour un rôle"
// @Success 200 {object} models.UpdateJobResponse "Rôle mis à jour avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/{id} [put]
// @Security ApiKeyAuth
func (h *UpdateJobHandler) Handle(c *gin.Context) {
	jobId := c.Param("id")
	if jobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	job, err := h.useCase.Execute(c.Request.Context(), entity.JobId(jobId), request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}
