package get_job_for_user

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetJobForUserHandler struct {
	useCase *GetJobForUserUseCase
}

func NewGetJobForUserHandler(useCase *GetJobForUserUseCase) *GetJobForUserHandler {
	return &GetJobForUserHandler{useCase: useCase}
}

// Handle récupère les rôles de travail attribués à un utilisateur
// @Summary Obtenir les rôles de travail d'un utilisateur
// @Description Récupère la liste des rôles de travail avec leurs permissions pour un utilisateur spécifique
// @Tags job
// @Accept json
// @Produce json
// @Param user_id path string true "ID de l'utilisateur"
// @Success 200 {object} get_job_for_user.JobResponse "Liste des rôles de travail de l'utilisateur"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/user/{user_id} [get]
// @Security ApiKeyAuth
func (h *GetJobForUserHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	jobs, err := h.useCase.Execute(c.Request.Context(), user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var jobResponses []JobResponse
	for _, job := range jobs {
		jobResponses = append(jobResponses, JobResponse{
			Id:          string(job.Id),
			Name:        job.Name,
			Permissions: int(job.Permissions),
			IsAssigned:  job.IsAssigned,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobResponses,
	})
}

type JobResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Permissions int    `json:"permissions"`
	IsAssigned  bool   `json:"is_assigned"`
}
