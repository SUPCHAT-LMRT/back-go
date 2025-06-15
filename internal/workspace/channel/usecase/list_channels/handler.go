package list_channels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ListChannelsHandler struct {
	useCase *ListChannelsUseCase
}

func NewListChannelsHandler(useCase *ListChannelsUseCase) *ListChannelsHandler {
	return &ListChannelsHandler{useCase: useCase}
}

// TODO: filter out the channels that the user is not a member of

// Handle récupère la liste des canaux d'un espace de travail
// @Summary Liste des canaux
// @Description Retourne tous les canaux publics d'un espace de travail
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Success 200 {array} ListChannelResponse "Liste des canaux"
// @Failure 400 {object} map[string]string "ID de workspace manquant"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des canaux"
// @Router /api/workspaces/{workspace_id}/channels [get]
// @Security ApiKeyAuth
func (h *ListChannelsHandler) Handle(c *gin.Context) {
	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	channels, err := h.useCase.Execute(c, entity.WorkspaceId(workspaceId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ListChannelResponse, len(channels))
	for i, channel := range channels {
		response[i] = ListChannelResponse{
			Id:    string(channel.Id),
			Name:  channel.Name,
			Topic: channel.Topic,
		}
	}

	c.JSON(http.StatusOK, response)
}

type ListChannelResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
