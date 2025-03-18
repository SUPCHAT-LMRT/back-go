package login_oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"net/http"
	"os"
)

type OAuthHandler struct {
	loginUseCase *OAuthUseCase
}

func NewOAuthHandler(loginUseCase *OAuthUseCase) *OAuthHandler {
	return &OAuthHandler{loginUseCase: loginUseCase}
}

// Démarrer l'authentification OAuth avec Google ou Facebook
func (h *OAuthHandler) AuthProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "facebook" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback après la connexion OAuth
func (h *OAuthHandler) AuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "facebook" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Vérifier ou enregistrer l'utilisateur en base de données
	response, err := h.loginUseCase.HandleOAuthLogin(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Stocker les tokens dans des cookies
	c.SetCookie("accessToken", response.AccessToken, int(response.AccessTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)
	c.SetCookie("refreshToken", response.RefreshToken, int(response.RefreshTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)

	c.JSON(http.StatusOK, h.formatUserResponse(response.User))
}

// Formater la réponse utilisateur
func (h *OAuthHandler) formatUserResponse(user *entity.User) gin.H {
	return gin.H{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		//"provider":   user.Provider,
	}
}
