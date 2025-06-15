package permissions

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
)

type CheckUserPermissionsHandler struct {
	checkPermissionUseCase *CheckPermissionJobUseCase
}

func NewCheckUserPermissionsHandler(
	useCase *CheckPermissionJobUseCase,
) *CheckUserPermissionsHandler {
	return &CheckUserPermissionsHandler{checkPermissionUseCase: useCase}
}

// Handle vérifie si un utilisateur possède les permissions demandées
// @Summary Vérifier les permissions d'un utilisateur
// @Description Vérifie si un utilisateur spécifié possède les permissions demandées
// @Tags job
// @Accept json
// @Produce json
// @Param user_id path string true "ID de l'utilisateur"
// @Param request body models.CheckPermissionsRequest true "Paramètres de vérification des permissions"
// @Success 200 {object} models.CheckPermissionsResponse "Résultat de la vérification"
// @Failure 400 {object} map[string]string "Erreur de paramètre"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/job/check-permissions/{user_id} [post]
// @Security ApiKeyAuth
func (h *CheckUserPermissionsHandler) Handle(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var request struct {
		Permissions uint64 `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hasPermission, err := h.checkPermissionUseCase.Execute(c, user_entity.UserId(userId), request.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasPermission": hasPermission})
}
