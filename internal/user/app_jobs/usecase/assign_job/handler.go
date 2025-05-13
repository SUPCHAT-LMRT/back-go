package assign_job

import (
	"github.com/gin-gonic/gin"
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

func (h *AssignJobHandler) Handle(c *gin.Context) {
	var request struct {
		JobId  string `json:"jobId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID and User ID are required"})
		return
	}

	err := h.useCase.Execute(c.Request.Context(), entity.JobId(request.JobId), user_entity.UserId(request.UserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job assigned successfully"})
}
