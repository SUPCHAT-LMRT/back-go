package check_permissions

import (
	"github.com/gin-gonic/gin"
	_ "github.com/supchat-lmrt/back-go/internal/models" // Import pour que Swagger trouve les modèles
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/permissions"
	"net/http"
)

type CheckPermissionsHandler struct {
	checkPermissionUseCase *permissions.CheckPermissionUseCase
}

func NewCheckPermissionsHandler(
	checkPermissionUseCase *permissions.CheckPermissionUseCase,
) *CheckPermissionsHandler {
	return &CheckPermissionsHandler{checkPermissionUseCase: checkPermissionUseCase}
}

// Handle vérifie si l'utilisateur possède des permissions spécifiques dans l'espace de travail
// @Summary Vérification de permissions
// @Description Vérifie si l'utilisateur courant possède un ensemble de permissions dans l'espace de travail
// @Tags workspace,permissions
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param request body models.CheckPermissionsRequest true "Permissions à vérifier"
// @Success 200 {object} models.CheckPermissionsResponse "Résultat de la vérification des permissions"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant ou requête invalide"
// @Failure 403 {object} map[string]string "Utilisateur non membre de l'espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la vérification des permissions"
// @Router /api/workspaces/{workspace_id}/permissions/check [post]
// @Security ApiKeyAuth
func (h *CheckPermissionsHandler) Handle(c *gin.Context) {
	workspaceMember := c.MustGet("workspace_member").(*workspace_entity.WorkspaceMember) //nolint:revive
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Permissions uint64 `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hasPermission, err := h.checkPermissionUseCase.Execute(
		c,
		workspaceMember.Id,
		entity2.WorkspaceId(workspaceId),
		request.Permissions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasPermission": hasPermission})
}
