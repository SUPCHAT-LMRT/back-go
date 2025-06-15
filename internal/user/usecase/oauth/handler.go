package oauth

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
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
	LinkOAuthUseCase     *LinkOAuthUseCase
}

type RegisterOAuthHandler struct {
	deps RegisterOAuthHandlerDeps
}

func NewRegisterOAuthHandler(deps RegisterOAuthHandlerDeps) *RegisterOAuthHandler {
	return &RegisterOAuthHandler{deps: deps}
}

// Provider starts the OAuth authentication with configured providers
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

func (h *RegisterOAuthHandler) LinkProvider(c *gin.Context) {
	user := c.MustGet("user").(*user_entity.User)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	q.Add("state", user.Id.String())
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *RegisterOAuthHandler) LinkCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	userId := c.Query("state")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId in state"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ici on appelle la use case pour lier l'account OAuth à l'utilisateur existant
	err = h.deps.LinkOAuthUseCase.Execute(c, userId, provider, oauthUser.UserID, oauthUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("FRONT_ACCOUNT_SETTINGS_URL")) // Retour à la page profil/settings
}
