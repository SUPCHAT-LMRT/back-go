package get_job_for_user

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}
