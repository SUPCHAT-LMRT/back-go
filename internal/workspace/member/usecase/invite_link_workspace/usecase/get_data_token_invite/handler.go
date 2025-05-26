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
