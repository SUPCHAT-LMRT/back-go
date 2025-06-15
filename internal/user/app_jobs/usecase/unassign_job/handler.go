package unassign_job

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type UnassignJobHandler struct {
	useCase *UnassignJobUseCase
}

func NewUnassignJobHandler(useCase *UnassignJobUseCase) *UnassignJobHandler {
	return &UnassignJobHandler{useCase: useCase}
}

// Handle désassigne un rôle de travail d'un utilisateur
// @Summary Désassigner un rôle de travail
// @Description Retire un rôle de travail spécifié d'un utilisateur
// @Tags job
// @Accept json
// @Produce json
// @Param request body models.UnassignJobRequest true "Informations pour retirer un rôle"
// @Success 200 {object} models.UnassignJobResponse "Opération réussie"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/unassign [post]
// @Security ApiKeyAuth
func (h *UnassignJobHandler) Handle(c *gin.Context) {
	var request struct {
		JobId  string `json:"jobId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID and User ID are required"})
		return
	}

	err := h.useCase.Execute(
		c.Request.Context(),
		entity.JobId(request.JobId),
		user_entity.UserId(request.UserId),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job unassigned successfully"})
}
