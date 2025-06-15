package delete_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type DeleteUserHandler struct {
	useCase *DeleteUserUseCase
}

func NewDeleteUserHandler(useCase *DeleteUserUseCase) *DeleteUserHandler {
	return &DeleteUserHandler{useCase: useCase}
}

// Handle supprime un utilisateur du système
// @Summary Supprimer un utilisateur
// @Description Supprime définitivement un utilisateur du système
// @Tags account
// @Accept json
// @Produce json
// @Param userId path string true "ID de l'utilisateur à supprimer"
// @Success 200 {object} map[string]string "Utilisateur supprimé avec succès"
// @Failure 400 {object} map[string]string "ID utilisateur manquant"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/auth/delete/{userId} [delete]
// @Security ApiKeyAuth
func (h *DeleteUserHandler) Handle(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	err := h.useCase.Execute(c.Request.Context(), user_entity.UserId(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
