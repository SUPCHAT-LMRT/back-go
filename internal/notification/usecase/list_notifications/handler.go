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

// Handle récupère la liste des notifications de l'utilisateur
// @Summary Lister les notifications
// @Description Récupère toutes les notifications de l'utilisateur connecté
// @Tags notification
// @Accept json
// @Produce json
// @Success 200 {array} entity.Notification "Liste des notifications"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/notifications [get]
// @Security ApiKeyAuth
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
