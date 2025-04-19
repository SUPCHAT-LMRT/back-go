package oauth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/register"
	uberdig "go.uber.org/dig"
	"net/http"
	"os"
	"strings"
)

type RegisterOAuthHandlerDeps struct {
	uberdig.In
	LoginOAuthUseCase    *LoginOAuthUseCase
	RegisterOAuthUseCase *RegisterOAuthUseCase
	Logger               logger.Logger
}

type RegisterOAuthHandler struct {
	deps RegisterOAuthHandlerDeps
}

func NewRegisterOAuthHandler(deps RegisterOAuthHandlerDeps) *RegisterOAuthHandler {
	deps.Logger = deps.Logger.With().Str("handler", "RegisterOAuthHandler").Logger()
	return &RegisterOAuthHandler{deps: deps}
}

// Provider starts the OAuth authentication with configured providers
func (h *RegisterOAuthHandler) Provider(c *gin.Context) {
	provider := c.Param("provider")
	inviteToken := c.Query("token")
	if provider != "google" && provider != "facebook" {
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
	if provider != "google" && provider != "facebook" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	inviteToken := c.Query("state")

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(inviteToken) != 36 {
		// If no invite token is provided (not a uuid), redirect to the home page
		response, err := h.deps.LoginOAuthUseCase.Execute(c, user.Email)
		if err != nil {
			if errors.Is(err, user_repository.UserNotFoundErr) {
				c.Redirect(http.StatusFound, os.Getenv("FRONT_ACCOUNT_LOGIN_URL")+"?error=Utilisateur non trouvé")
				err = gothic.Logout(c.Writer, c.Request)
				if err != nil {
					h.deps.Logger.Error().Err(err).Msg("Error while logging out after logging in")
					return
				}
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Stocker les tokens dans des cookies
		c.SetCookie("accessToken", response.AccessToken, int(response.AccessTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)
		c.SetCookie("refreshToken", response.RefreshToken, int(response.RefreshTokenLifespan.Seconds()), "/", os.Getenv("DOMAIN"), false, true)

		c.Redirect(http.StatusFound, os.Getenv("FRONT_URL"))
		return
	}

	// If an invitation token is provided, redirect to the login page
	err = h.deps.RegisterOAuthUseCase.Execute(c, inviteToken, user.Email)
	if err != nil {
		if errors.Is(err, register.UserAlreadyExistsErr) {
			redirectUrl := os.Getenv("FRONT_ACCOUNT_REGISTER_URL")
			redirectUrl = strings.Replace(redirectUrl, "{token}", inviteToken, -1)
			c.Redirect(http.StatusFound, redirectUrl+"?error=Un utilisateur avec cette adresse email existe déjà")
			err = gothic.Logout(c.Writer, c.Request)
			if err != nil {
				h.deps.Logger.Error().Err(err).Msg("Error while logging out after registering")
				return
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("FRONT_ACCOUNT_LOGIN_URL")+"?success=Connecté avec succès")
}
