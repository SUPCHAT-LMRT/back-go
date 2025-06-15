package create_job

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	"net/http"
)

type CreateJobHandler struct {
	useCase *CreateJobUseCase
}

func NewCreateJobHandler(useCase *CreateJobUseCase) *CreateJobHandler {
	return &CreateJobHandler{useCase: useCase}
}

// Handle crée un nouveau rôle de travail dans le système
// @Summary Créer un rôle de travail
// @Description Crée un nouveau rôle de travail avec le nom fourni
// @Tags job
// @Accept json
// @Produce json
// @Param request body models.CreateJobRequest true "Informations du rôle à créer"
// @Success 200 {object} models.CreateJobResponse "Rôle créé avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 409 {object} map[string]string "Conflit - Le rôle existe déjà"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job [post]
// @Security ApiKeyAuth
func (h *CreateJobHandler) Handle(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	job, err := h.useCase.Execute(c, request.Name)
	if err != nil {
		if err.Error() == fmt.Sprintf("a job with the name '%s' already exists", request.Name) {
			c.JSON(
				http.StatusConflict,
				gin.H{"error": fmt.Sprintf("Job with name '%s' already exists", request.Name)},
			)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	response := createJobResponse{
		Id:   string(job.Id),
		Name: job.Name,
	}

	c.JSON(http.StatusOK, response)
}

type createJobResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
