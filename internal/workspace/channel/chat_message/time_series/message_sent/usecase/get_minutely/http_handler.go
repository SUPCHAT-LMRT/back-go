package get_minutely

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetMinutelyMessageSentHandler struct {
	useCase *GetMinutelyMessageSentUseCase
}

func NewGetMinutelyMessageSentHandler(
	useCase *GetMinutelyMessageSentUseCase,
) *GetMinutelyMessageSentHandler {
	return &GetMinutelyMessageSentHandler{useCase: useCase}
}

func (h GetMinutelyMessageSentHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	messageSents, err := h.useCase.Execute(
		c,
		entity.WorkspaceId(workspaceId),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messageSentsResponse := make([]*MessageSentResponse, len(messageSents))
	for i, sent := range messageSents {
		messageSentsResponse[i] = &MessageSentResponse{
			SentAt: sent.SentAt,
			Count:  sent.Count,
		}
	}

	c.JSON(http.StatusOK, messageSentsResponse)
}

type MessageSentResponse struct {
	SentAt time.Time `json:"sentAt"`
	Count  uint      `json:"count"`
}
