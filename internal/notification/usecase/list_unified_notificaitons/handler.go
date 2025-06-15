package list_unified_notifications

//
//import (
//	"net/http"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
//)
//
//type ListUnifiedNotificationsHandler struct {
//	useCase *ListUnifiedNotificationsUseCase
//}
//
//func NewListUnifiedNotificationsHandler(useCase *ListUnifiedNotificationsUseCase) *ListUnifiedNotificationsHandler {
//	return &ListUnifiedNotificationsHandler{useCase: useCase}
//}
//
//func (h *ListUnifiedNotificationsHandler) Handle(c *gin.Context) {
//	// Récupérer l'utilisateur depuis le contexte (middleware d'auth)
//	userIdValue, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
//		return
//	}
//
//	userId, ok := userIdValue.(user_entity.UserId)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID utilisateur invalide"})
//		return
//	}
//
//	// Paramètres de pagination
//	limitStr := c.DefaultQuery("limit", "20")
//	offsetStr := c.DefaultQuery("offset", "0")
//
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil || limit < 1 || limit > 100 {
//		limit = 20
//	}
//
//	offset, err := strconv.Atoi(offsetStr)
//	if err != nil || offset < 0 {
//		offset = 0
//	}
//
//	req := ListUnifiedNotificationsRequest{
//		UserId: userId,
//		Limit:  limit,
//		Offset: offset,
//	}
//
//	response, err := h.useCase.Execute(c.Request.Context(), req)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des notifications"})
//		return
//	}
//
//	c.JSON(http.StatusOK, response)
//}
