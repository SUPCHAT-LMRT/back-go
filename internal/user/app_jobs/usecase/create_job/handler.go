package create_job

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateJobHandler struct {
	useCase *CreateJobUseCase
}

func NewCreateJobHandler(useCase *CreateJobUseCase) *CreateJobHandler {
	return &CreateJobHandler{useCase: useCase}
}

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
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Job with name '%s' already exists", request.Name)})
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
