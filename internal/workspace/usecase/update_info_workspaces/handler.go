package update_info_workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateInfoWorkspacesHandler struct {
	useCase *UpdateInfoWorkspacesUseCase
}

func NewUpdateInfoWorkspacesHandler(
	useCase *UpdateInfoWorkspacesUseCase,
) *UpdateInfoWorkspacesHandler {
	return &UpdateInfoWorkspacesHandler{useCase: useCase}
}

// Handle met à jour les informations d'un espace de travail
// @Summary Mise à jour des informations de l'espace de travail
// @Description Modifie le nom et le sujet d'un espace de travail
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param request body object true "Informations à mettre à jour"
// @Success 200 {object} nil "Informations mises à jour avec succès"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant ou requête invalide"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la mise à jour des informations"
// @Router /api/workspaces/{workspace_id} [put]
// @Security ApiKeyAuth
func (h *UpdateInfoWorkspacesHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	var request struct {
		Name  string `json:"name"`
		Topic string `json:"topic"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request"})
		return
	}

	err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), request.Name, request.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
