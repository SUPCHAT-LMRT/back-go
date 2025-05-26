package delete_job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
)

type DeleteJobHandler struct {
	useCase *DeleteJobUseCase
}

func NewDeleteJobHandler(useCase *DeleteJobUseCase) *DeleteJobHandler {
	return &DeleteJobHandler{useCase: useCase}
}

func (h *DeleteJobHandler) Handle(c *gin.Context) {
	jobId := c.Param("id")
	if jobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	err := h.useCase.Execute(c.Request.Context(), entity.JobId(jobId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}
