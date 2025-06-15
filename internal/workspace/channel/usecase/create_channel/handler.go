package create_channel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type CreateChannelHandler struct {
	useCase *CreateChannelUseCase
}

func NewCreateChannelHandler(useCase *CreateChannelUseCase) *CreateChannelHandler {
	return &CreateChannelHandler{useCase: useCase}
}

// Handle crée un nouveau canal dans un espace de travail
// @Summary Création d'un canal
// @Description Crée un nouveau canal textuel dans un espace de travail spécifique
// @Tags workspace,channel
// @Accept json
// @Produce json
// @Param workspace_id path string true "ID de l'espace de travail"
// @Param channel body CreateChannelRequest true "Informations du canal à créer"
// @Success 201 {string} string "Canal créé avec succès"
// @Failure 400 {object} map[string]string "Paramètres manquants ou invalides"
// @Failure 403 {object} map[string]string "Permissions insuffisantes pour créer un canal"
// @Failure 500 {object} map[string]string "Erreur lors de la création du canal"
// @Router /api/workspaces/{workspace_id}/channels [post]
// @Security ApiKeyAuth
func (h *CreateChannelHandler) Handle(c *gin.Context) {
	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspaceId := c.Param("workspace_id")
	if workspaceId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id is required"})
		return
	}

	if req.Members == nil {
		req.Members = make([]string, 0)
	}

	err := h.useCase.Execute(c, &channel_entity.Channel{
		Name:        req.Name,
		Topic:       req.Topic,
		WorkspaceId: entity.WorkspaceId(workspaceId),
		Kind:        channel_entity.ChannelKindText,
		IsPrivate:   req.IsPrivate,
		Members:     req.Members,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

type CreateChannelRequest struct {
	Name      string   `json:"name"      binding:"required,min=1,max=100"`
	Topic     string   `json:"topic"`
	IsPrivate bool     `json:"isPrivate"`
	Members   []string `json:"members"`
}
