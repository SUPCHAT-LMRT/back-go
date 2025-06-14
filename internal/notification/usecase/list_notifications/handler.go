package list_notifications

import (
	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type ListNotificationsHandler struct {
	ListNotificationsUseCase *ListNotificationsUseCase
}

func NewListNotificationsHandler(useCase *ListNotificationsUseCase) *ListNotificationsHandler {
	return &ListNotificationsHandler{
		ListNotificationsUseCase: useCase,
	}
}

func (h *ListNotificationsHandler) Handle(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	notifications, err := h.ListNotificationsUseCase.Execute(c.Request.Context(), user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}
