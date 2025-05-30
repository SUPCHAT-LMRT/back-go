package get_poll_by_id

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetPollByIdHandler struct {
	usecase *GetPollByIdUseCase
}

func NewGetPollByIdHandler(usecase *GetPollByIdUseCase) *GetPollByIdHandler {
	return &GetPollByIdHandler{usecase: usecase}
}

func (h *GetPollByIdHandler) Handle(c *gin.Context) {
	pollId := c.Param("poll_id")
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poll_id is required"})
		return
	}

	poll, err := h.usecase.Execute(c, pollId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, poll)
}
