package list_private_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type GetPrivateChannelsHandler struct {
	useCase *GetPrivateChannelsUseCase
}

func NewGetPrivateChannelsHandler(useCase *GetPrivateChannelsUseCase) *GetPrivateChannelsHandler {
	return &GetPrivateChannelsHandler{useCase: useCase}
}

// Handle récupère la liste des canaux privés accessibles par l'utilisateur
// @Summary Liste des canaux privés
// @Description Retourne tous les canaux privés auxquels l'utilisateur a accès dans un espace de travail
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {array} GetPrivateChannelResponse "Liste des canaux privés accessibles"
// @Failure 400 {object} map[string]string "ID de workspace manquant"
// @Failure 401 {object} map[string]string "Utilisateur non authentifié ou non membre du workspace"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des canaux"
// @Router /api/workspaces/{workspace_id}/channels/private [get]
// @Security ApiKeyAuth
func (h *GetPrivateChannelsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	workspaceMember, ok := c.MustGet("workspace_member").(*workspace_member_entity.WorkspaceMember)
	if !ok || workspaceMember == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id and user_id are required"})
		return
	}

	channels, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId), workspaceMember.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]GetPrivateChannelResponse, len(channels))
	for i, channel := range channels {
		response[i] = GetPrivateChannelResponse{
			Id:    string(channel.Id),
			Name:  channel.Name,
			Topic: channel.Topic,
		}
	}

	c.JSON(http.StatusOK, response)
}

type GetPrivateChannelResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
