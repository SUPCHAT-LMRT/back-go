package get_channel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type GetChannelHandler struct {
	useCase *GetChannelUseCase
}

func NewGetChannelHandler(useCase *GetChannelUseCase) *GetChannelHandler {
	return &GetChannelHandler{useCase: useCase}
}

// Handle récupère les informations d'un canal spécifique
// @Summary Détails d'un canal
// @Description Récupère les informations détaillées d'un canal dans un espace de travail
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channel_id path string true "ID du canal"
// @Success 200 {object} ChannelResponse "Informations détaillées du canal"
// @Failure 404 {object} map[string]string "Canal non trouvé"
// @Failure 500 {object} map[string]string "Erreur lors de la récupération des informations du canal"
// @Router /api/workspaces/{workspace_id}/channels/{channel_id} [get]
// @Security ApiKeyAuth
func (h *GetChannelHandler) Handle(c *gin.Context) {
	channelId := c.Param("channel_id")
	channel, err := h.useCase.Execute(c.Request.Context(), entity.ChannelId(channelId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         channel.Id,
		"name":       channel.Name,
		"topic":      channel.Topic,
		"kind":       channel.Kind,
		"isPrivate":  channel.IsPrivate,
		"members":    channel.Members,
		"workspace":  channel.WorkspaceId,
		"created_at": channel.CreatedAt,
		"updated_at": channel.UpdatedAt,
		"index":      channel.Index,
	})
}

type ChannelResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	WorkspaceId string `json:"workspaceId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
