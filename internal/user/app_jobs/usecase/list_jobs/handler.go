package list_jobs

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	"net/http"
	"sort"
	"strconv"
)

type ListJobsHandler struct {
	useCase *ListJobsUseCase
}

func NewListJobsHandler(useCase *ListJobsUseCase) *ListJobsHandler {
	return &ListJobsHandler{useCase: useCase}
}

//nolint:revive

// Handle récupère la liste de tous les rôles de travail dans le système
// @Summary Lister les rôles de travail
// @Description Récupère la liste complète des rôles de travail avec leurs ID, noms et permissions
// @Tags job
// @Accept json
// @Produce json
// @Success 200 {object} models.ListJobsResponse "Liste des rôles de travail"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job [get]
// @Security ApiKeyAuth
func (h *ListJobsHandler) Handle(c *gin.Context) {
	jobs, err := h.useCase.Execute(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
		return
	}

	sort.Slice(jobs, func(i, j int) bool {
		if jobs[i].Name == "Admin" {
			return true
		}
		if jobs[j].Name == "Admin" {
			return false
		}
		if jobs[i].Name == "Manager" {
			return true
		}
		if jobs[j].Name == "Manager" {
			return false
		}
		return jobs[i].Name < jobs[j].Name
	})

	var jobResponses []jobResponse
	for _, job := range jobs {
		jobResponses = append(jobResponses, jobResponse{
			ID:         string(job.Id),
			Name:       job.Name,
			Permission: strconv.FormatUint(job.Permissions, 10),
		})
	}

	c.JSON(http.StatusOK, listJobsResponse{
		Jobs: jobResponses,
	})
}

type listJobsResponse struct {
	Jobs []jobResponse `json:"jobs"`
}

type jobResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
}
