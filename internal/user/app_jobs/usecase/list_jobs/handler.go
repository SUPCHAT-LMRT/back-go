package list_jobs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListJobsHandler struct {
	useCase *ListJobsUseCase
}

func NewListJobsHandler(useCase *ListJobsUseCase) *ListJobsHandler {
	return &ListJobsHandler{useCase: useCase}
}

func (h *ListJobsHandler) Handle(c *gin.Context) {
	jobs, err := h.useCase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}
