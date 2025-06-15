package generate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateInviteLinkRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"  binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
}

type CreateInviteLinkHandler struct {
	usecase *InviteLinkUseCase
}

func NewCreateInviteLinkHandler(usecase *InviteLinkUseCase) *CreateInviteLinkHandler {
	return &CreateInviteLinkHandler{usecase: usecase}
}

// Handle génère un lien d'invitation pour un nouvel utilisateur
// @Summary Créer un lien d'invitation
// @Description Crée un lien d'invitation pour permettre à un nouvel utilisateur de s'inscrire
// @Tags account
// @Accept json
// @Produce plain
// @Param request body generate.CreateInviteLinkRequest true "Informations de l'utilisateur invité"
// @Success 200 {string} string "URL du lien d'invitation"
// @Failure 400 {object} map[string]string "Paramètres de requête invalides"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Failure 500 {object} map[string]string "Erreur interne du serveur"
// @Router /api/account/invite-link [post]
// @Security ApiKeyAuth
func (h *CreateInviteLinkHandler) Handle(c *gin.Context) {
	var req CreateInviteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviteLink, err := h.usecase.CreateInviteLink(c, req.FirstName, req.LastName, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, inviteLink)
}
