package mark_as_read

import (
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"net/http"
)

type MarkAsReadHandler struct {
	MarkAsReadUseCase *MarkAsReadUseCase
}

func NewMarkAsReadHandler(useCase *MarkAsReadUseCase) *MarkAsReadHandler {
	return &MarkAsReadHandler{
		MarkAsReadUseCase: useCase,
	}
}

// Handle marque une notification comme lue
// @Summary Marquer une notification comme lue
// @Description Change l'état d'une notification spécifique pour la marquer comme lue
// @Tags notification
// @Accept json
// @Produce json
// @Param id path string true "ID de la notification"
// @Success 204 "Notification marquée comme lue avec succès"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/notifications/{id}/read [patch]
// @Security ApiKeyAuth
func (h *MarkAsReadHandler) Handle(c *gin.Context) {
	notificationId := c.Param("id")
	if notificationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	err := h.MarkAsReadUseCase.Execute(c.Request.Context(), entity.NotificationId(notificationId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
