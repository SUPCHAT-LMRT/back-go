package get_workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetWorkspaceHandler struct {
	useCase *GetWorkspaceUseCase
}

func NewGetWorkspaceHandler(useCase *GetWorkspaceUseCase) *GetWorkspaceHandler {
	return &GetWorkspaceHandler{useCase: useCase}
}

// Handle récupère les informations d'un espace de travail
// @Summary Détails d'un espace de travail
// @Description Récupère les informations de base d'un espace de travail spécifique
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {object} map[string]string "Informations de l'espace de travail"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération de l'espace de travail"
// @Router /api/workspaces/{workspace_id} [get]
// @Security ApiKeyAuth
func (h *GetWorkspaceHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	workspace, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    workspace.Id,
		"name":  workspace.Name,
		"topic": workspace.Topic,
		"type":  workspace.Type,
	})
}
