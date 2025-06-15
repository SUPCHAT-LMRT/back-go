package get_workspace_details

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type GetWorkspaceDetailsHandler struct {
	useCase *GetWorkspaceDetailsUseCase
}

func NewGetWorkspaceDetailsHandler(
	useCase *GetWorkspaceDetailsUseCase,
) *GetWorkspaceDetailsHandler {
	return &GetWorkspaceDetailsHandler{useCase: useCase}
}

// Handle récupère les détails complets d'un espace de travail
// @Summary Détails complets d'un espace de travail
// @Description Récupère les informations détaillées d'un espace de travail spécifique avec statistiques
// @Tags workspace
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {object} get_workspace_details.WorkspaceDetailsResponse "Détails complets de l'espace de travail"
// @Failure 400 {object} map[string]string "ID de l'espace de travail manquant"
// @Failure 403 {object} map[string]string "Utilisateur non autorisé dans cet espace de travail"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des détails"
// @Router /api/workspaces/{workspace_id}/details [get]
// @Security ApiKeyAuth
func (h *GetWorkspaceDetailsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	workspaceDetails, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, WorkspaceDetailsResponse{
		Id:            workspaceDetails.Id.String(),
		Name:          workspaceDetails.Name,
		Topic:         workspaceDetails.Topic,
		Type:          string(workspaceDetails.Type),
		MembersCount:  workspaceDetails.MembersCount,
		ChannelsCount: workspaceDetails.ChannelsCount,
		MessagesCount: workspaceDetails.MessagesCount,
	})
}

type WorkspaceDetailsResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Topic         string `json:"topic"`
	Type          string `json:"type"`
	MembersCount  uint   `json:"membersCount"`
	ChannelsCount uint   `json:"channelsCount"`
	MessagesCount uint   `json:"messagesCount"`
}
