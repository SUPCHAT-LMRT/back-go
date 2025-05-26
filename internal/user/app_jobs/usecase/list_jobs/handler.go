package list_jobs

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

type ListJobsHandler struct {
	useCase *ListJobsUseCase
}

func NewListJobsHandler(useCase *ListJobsUseCase) *ListJobsHandler {
	return &ListJobsHandler{useCase: useCase}
}

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
