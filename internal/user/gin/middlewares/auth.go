package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
)

type AuthMiddleware struct {
	getUserByEmail *get_by_email.GetUserByEmailUseCase
	tokenStrategy  token.TokenStrategy
}

func NewAuthMiddleware(
	getUserByEmail *get_by_email.GetUserByEmailUseCase,
	tokenStrategy token.TokenStrategy,
) *AuthMiddleware {
	return &AuthMiddleware{getUserByEmail: getUserByEmail, tokenStrategy: tokenStrategy}
}

//nolint:revive
func (a *AuthMiddleware) Execute(c *gin.Context) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil || accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token required"})
		return
	}

	accessTokenClaims, err := a.tokenStrategy.ValidateAccessToken(accessToken)
	if err != nil {
		if errors.Is(err, token.ErrUnableToValidateAccessToken) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"message": "unable to validate access token"},
			)
			return
		}
		if errors.Is(err, token.ErrAccessTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "access token expired"})
			return
		}
		if errors.Is(err, token.ErrInvalidTokenType) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token type"})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid access token"})
		return
	}

	user, err := a.getUserByEmail.Execute(c, accessTokenClaims["email"].(string))
	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"message": "Unable to find user relative to this account!"},
			)
			return
		}
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Internal server error"},
		)
		return
	}

	c.Set("user", user)
}
