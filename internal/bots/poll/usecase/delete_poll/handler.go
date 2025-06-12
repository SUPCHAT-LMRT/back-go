package delete_poll

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeletePollHandler struct {
	usecase *DeletePollUseCase
}

func NewDeletePollHandler(usecase *DeletePollUseCase) *DeletePollHandler {
	return &DeletePollHandler{usecase: usecase}
}

func (h *DeletePollHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poll_id is required"})
		return
	}

	err := h.usecase.Execute(c, pollId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poll deleted successfully"})
}
