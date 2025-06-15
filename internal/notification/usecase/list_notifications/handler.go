package list_notifications

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
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
	user := c.MustGet("user").(*user_entity.User) //nolint:revive

	notifications, err := h.ListNotificationsUseCase.Execute(c.Request.Context(), user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(notifications) == 0 {
		notifications = []*entity.Notification{}
	}
	c.JSON(http.StatusOK, notifications)
}
