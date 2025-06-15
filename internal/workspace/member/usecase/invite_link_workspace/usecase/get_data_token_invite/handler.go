package get_data_token_invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
)

type GetInviteLinkWorkspaceDataHandler struct {
	usecase             *GetInviteLinkDataUseCase
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase
}

func NewGetInviteLinkWorkspaceDataHandler(
	usecase *GetInviteLinkDataUseCase,
	getWorkspaceUseCase *get_workspace.GetWorkspaceUseCase,
) *GetInviteLinkWorkspaceDataHandler {
	return &GetInviteLinkWorkspaceDataHandler{
		usecase:             usecase,
		getWorkspaceUseCase: getWorkspaceUseCase,
	}
}

// Handle récupère les informations d'un lien d'invitation à un espace de travail
// @Summary Informations d'un lien d'invitation
// @Description Récupère les informations associées à un lien d'invitation à un espace de travail
// @Tags workspace,invite
// @Accept json
// @Produce json
// @Param token path string true "Token d'invitation à vérifier"
// @Success 200 {object} InviteLinkDataResponse "Informations sur l'espace de travail lié à l'invitation"
// @Failure 400 {object} map[string]string "Token manquant"
// @Failure 404 {object} map[string]string "Token d'invitation invalide ou expiré"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des informations"
// @Router /api/workspace-invite-link/{token} [get]
func (h *GetInviteLinkWorkspaceDataHandler) Handle(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	inviteLink, err := h.usecase.GetInviteLinkData(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workspace, err := h.getWorkspaceUseCase.Execute(c, inviteLink.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, InviteLinkDataResponse{
		WorkspaceId:   inviteLink.WorkspaceId,
		WorkspaceName: workspace.Name,
	})
}

type InviteLinkDataResponse struct {
	WorkspaceId   entity.WorkspaceId `json:"workspaceId"`
	WorkspaceName string             `json:"workspaceName"`
}
