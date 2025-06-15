package assign_job

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type AssignJobHandler struct {
	useCase *AssignJobUseCase
}

func NewAssignJobHandler(useCase *AssignJobUseCase) *AssignJobHandler {
	return &AssignJobHandler{useCase: useCase}
}

// Handle attribue un rôle de travail à un utilisateur
// @Summary Assigner un rôle de travail
// @Description Attribue un rôle de travail spécifique à un utilisateur
// @Tags job
// @Accept json
// @Produce json
// @Param request body models.AssignJobRequest true "Informations d'attribution de rôle"
// @Success 200 {object} map[string]string "Rôle attribué avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/assign [post]
// @Security ApiKeyAuth
func (h *AssignJobHandler) Handle(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Job assigned successfully"})
}
