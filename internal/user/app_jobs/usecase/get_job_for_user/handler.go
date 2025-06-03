package get_job_for_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetJobForUserHandler struct {
	useCase *GetJobForUserUseCase
}

func NewGetJobForUserHandler(useCase *GetJobForUserUseCase) *GetJobForUserHandler {
	return &GetJobForUserHandler{useCase: useCase}
}

func (h *GetJobForUserHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	jobs, err := h.useCase.Execute(c.Request.Context(), userId)
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
