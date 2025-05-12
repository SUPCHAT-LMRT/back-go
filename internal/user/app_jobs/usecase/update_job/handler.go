package update_job

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"net/http"
)

type UpdateJobHandler struct {
	useCase *UpdateJobUseCase
}

func NewUpdateJobHandler(useCase *UpdateJobUseCase) *UpdateJobHandler {
	return &UpdateJobHandler{useCase: useCase}
}

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
