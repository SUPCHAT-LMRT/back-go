package oauth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	uberdig "go.uber.org/dig"
)

type RegisterOAuthHandlerDeps struct {
	uberdig.In
	LoginOAuthUseCase    *LoginOAuthUseCase
	RegisterOAuthUseCase *RegisterOAuthUseCase
}

type RegisterOAuthHandler struct {
	deps RegisterOAuthHandlerDeps
}

func NewRegisterOAuthHandler(deps RegisterOAuthHandlerDeps) *RegisterOAuthHandler {
	return &RegisterOAuthHandler{deps: deps}
}

// Provider starts the OAuth authentication with configured providers

// Provider démarre l'authentification OAuth avec les fournisseurs configurés
// @Summary Authentification OAuth
// @Description Initie le processus d'authentification OAuth avec un fournisseur spécifié
// @Tags account
// @Accept json
// @Produce json
// @Param provider path string true "Nom du fournisseur OAuth (google ou github)"
// @Param token query string false "Token d'invitation optionnel"
// @Success 302 {string} string "Redirection vers le fournisseur OAuth"
// @Failure 400 {object} map[string]string "Fournisseur OAuth invalide"
// @Router /api/account/auth/oauth/{provider} [get]
func (h *RegisterOAuthHandler) Provider(c *gin.Context) {
	provider := c.Param("provider")
	inviteToken := c.Query("token")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	q.Add("state", inviteToken)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback after OAuth authentication

// Callback traite la réponse après authentification OAuth
// @Summary Callback OAuth
// @Description Traite la réponse du fournisseur OAuth après authentification
// @Tags account
// @Accept json
// @Produce json
// @Param provider path string true "Nom du fournisseur OAuth (google ou github)"
// @Param state query string false "Token d'invitation utilisé lors de l'initialisation"
// @Success 302 {string} string "Redirection vers l'application ou la page de connexion"
// @Failure 400 {object} map[string]string "Fournisseur OAuth invalide"
// @Failure 500 {object} map[string]string "Erreur lors de l'authentification"
// @Router /api/account/auth/oauth/{provider}/callback [get]
func (h *RegisterOAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	inviteToken := c.Query("state")

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	oauthUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the state is a uuid (if so, it's an invite token)
	_, err = uuid.Parse(inviteToken)
	if err != nil {
		// If no invite token is provided, redirect to the home page
		response, err := h.deps.LoginOAuthUseCase.Execute(c, oauthUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Stocker les tokens dans des cookies
		c.SetCookie(
			"accessToken",
			response.AccessToken,
			int(response.AccessTokenLifespan.Seconds()),
			"/",
			os.Getenv("DOMAIN"),
			false,
			true,
		)
		c.SetCookie(
			"refreshToken",
			response.RefreshToken,
			int(response.RefreshTokenLifespan.Seconds()),
			"/",
			os.Getenv("DOMAIN"),
			false,
			true,
		)

		c.Redirect(http.StatusFound, os.Getenv("FRONT_URL"))
		return
	}

	// If an invite token is provided, redirect to the login page
	err = h.deps.RegisterOAuthUseCase.Execute(
		c,
		provider,
		oauthUser.UserID,
		oauthUser.Email,
		inviteToken,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("FRONT_ACCOUNT_LOGIN_URL"))
}
