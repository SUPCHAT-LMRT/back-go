package login_oauth

import (
	"context"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	"os"
	"time"
)

func init() {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLECREDID"), os.Getenv("GOOGLECREDSECRET"), os.Getenv("GOOGLECREDCALLBACKURL")),
		facebook.New(os.Getenv("FBCREDCREDID"), os.Getenv("FBCREDCREDSECRET"), os.Getenv("FBCREDCREDCALLBACKURL")),
	)
}

type OAuthUseCase struct {
	getUserByEmailUseCase *get_by_email.GetUserByEmailUseCase
	tokenStrategy         token.TokenStrategy
}

func NewOAuthUseCase(getUserByEmailUseCase *get_by_email.GetUserByEmailUseCase, tokenStrategy token.TokenStrategy) *OAuthUseCase {
	return &OAuthUseCase{getUserByEmailUseCase: getUserByEmailUseCase, tokenStrategy: tokenStrategy}
}

// Gestion de l'authentification via Google/Facebook
func (u *OAuthUseCase) HandleOAuthLogin(ctx context.Context, gothUser goth.User) (*OAuthResult, error) {
	// Vérifier si l'utilisateur existe déjà
	user, err := u.getUserByEmailUseCase.Execute(ctx, gothUser.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, le créer
		user = &entity.User{
			Email:      gothUser.Email,
			FirstName:  gothUser.FirstName,
			LastName:   gothUser.LastName,
			IsVerified: true,
		}
	}

	// Générer des tokens JWT
	return u.generateTokens(user, false)
}

// Génération des tokens d'accès et de rafraîchissement
func (u *OAuthUseCase) generateTokens(user *entity.User, rememberMe bool) (*OAuthResult, error) {
	accessToken, err := u.tokenStrategy.GenerateAccessToken(map[string]any{
		"email": user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := u.tokenStrategy.GenerateRefreshToken(map[string]any{
		"email": user.Email,
	}, rememberMe)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &OAuthResult{
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenLifespan:  u.tokenStrategy.GetAccessTokenLifespan(),
		RefreshTokenLifespan: u.tokenStrategy.GetRefreshTokenLifespan(rememberMe),
		User:                 user,
	}, nil
}

type OAuthResult struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenLifespan  time.Duration
	RefreshTokenLifespan time.Duration
	User                 *entity.User
}
