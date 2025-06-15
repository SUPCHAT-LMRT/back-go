package logout

import (
	"os"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct{}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

// Handle déconnecte l'utilisateur en supprimant ses cookies d'authentification
// @Summary Déconnexion utilisateur
// @Description Supprime les cookies d'authentification pour déconnecter l'utilisateur
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Déconnexion réussie"
// @Failure 401 {object} map[string]string "Non autorisé"
// @Router /api/account/auth/logout [post]
// @Security ApiKeyAuth
func (l LogoutHandler) Handle(c *gin.Context) {
	c.SetCookie("accessToken", "", -1, "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refreshToken", "", -1, "/", os.Getenv("DOMAIN"), false, true)
}
