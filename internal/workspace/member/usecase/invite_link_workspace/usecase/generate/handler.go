package generate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
)

type CreateInviteLinkHandler struct {
	useCase             *InviteLinkUseCase
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase
}

func NewCreateInviteLinkHandler(
	usecase *InviteLinkUseCase,
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase,
) *CreateInviteLinkHandler {
	return &CreateInviteLinkHandler{useCase: usecase, getWorkspaceUseCase: getWorkspaceUseCase}
}

type CreateInviteLinkRequest struct {
	WorkspaceId entity.WorkspaceId `json:"workspaceId"`
}

// Handle génère un lien d'invitation pour un espace de travail
// @Summary Création d'un lien d'invitation
// @Description Génère un lien d'invitation pour permettre à d'autres utilisateurs de rejoindre l'espace de travail
// @Tags workspace,invite
// @Accept json
// @Produce plain
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param body body CreateInviteLinkRequest true "Données de l'espace de travail"
// @Success 200 {string} string "Lien d'invitation généré"
// @Failure 400 {object} map[string]string "Requête invalide ou ID d'espace de travail manquant"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour créer un lien d'invitation"
// @Failure 404 {object} map[string]string "Espace de travail non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la génération du lien d'invitation"
// @Router /api/workspace-invite-link/create/{workspace_id} [post]
// @Security ApiKeyAuth
func (h *CreateInviteLinkHandler) Handle(c *gin.Context) {
	var request CreateInviteLinkRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.WorkspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WorkspaceId cannot be empty"})
		return
	}

	_, err := h.getWorkspaceUseCase.Execute(c, request.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	inviteLink, err := h.useCase.Execute(c, request.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, inviteLink)
}
