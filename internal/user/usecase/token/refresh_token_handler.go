package token

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	useCase *RefreshAccessTokenUseCase
}

func NewRefreshTokenHandler(useCase *RefreshAccessTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{useCase: useCase}
}

// Handle renouvelle le token d'accès à partir du token de rafraîchissement
// @Summary Rafraîchissement du token d'accès
// @Description Génère un nouveau token d'accès à partir du token de rafraîchissement présent dans les cookies
// @Tags account
// @Accept json
// @Produce json
// @Success 204 {string} string "Token d'accès renouvelé avec succès"
// @Failure 401 {object} map[string]string "Token de rafraîchissement manquant ou invalide"
// @Failure 500 {object} map[string]string "Erreur lors du rafraîchissement du token"
// @Router /api/account/auth/token/access/renew [post]
func (g *RefreshTokenHandler) Handle(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Refresh token is required.",
		})
		return
	}

	response, err := g.useCase.Execute(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to refresh the token.",
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie(
		"accessToken",
		response.AccessToken,
		int(response.AccessTokenLifespan.Seconds()),
		"/",
		os.Getenv("DOMAIN"),
		false,
		true,
	)

	c.Status(http.StatusNoContent)
}
