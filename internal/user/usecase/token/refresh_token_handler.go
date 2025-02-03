package token

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type RefreshTokenHandler struct {
	useCase *RefreshAccessTokenUseCase
}

func NewRefreshTokenHandler(useCase *RefreshAccessTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{useCase: useCase}
}

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

	c.SetCookie("accessToken", response.AccessToken, int(response.AccessTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)

	c.Status(http.StatusNoContent)
}
