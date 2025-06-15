package update_type_workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateTypeWorkspaceHandler struct {
	useCase *UpdateTypeWorkspaceUseCase
}

func NewUpdateTypeWorkspaceHandler(
	useCase *UpdateTypeWorkspaceUseCase,
) *UpdateTypeWorkspaceHandler {
	return &UpdateTypeWorkspaceHandler{useCase: useCase}
}

// Handle met à jour le type d'un espace de travail
// @Summary Mise à jour du type d'espace de travail
// @Description Modifie le type d'un espace de travail (public/privé)
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param request body object true "Nouveau type d'espace de travail"
// @Success 200 {object} nil "Type mis à jour avec succès"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant ou requête invalide"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la mise à jour du type"
// @Router /api/workspaces/{workspace_id}/type [put]
// @Security ApiKeyAuth
func (h *UpdateTypeWorkspaceHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(400, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request"})
		return
	}

	err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), entity.WorkspaceType(request.Type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
