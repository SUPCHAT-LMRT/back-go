package oauth

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/register"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	uberdig "go.uber.org/dig"
	"os"
	"time"
)

func init() {
	// TOOD Implement a shared store maybe ? What is exactly the store ?
	store := sessions.NewCookieStore([]byte("secret"))
	gothic.Store = store
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CRED_ID"), os.Getenv("GOOGLE_CRED_SECRET"), os.Getenv("GOOGLE_CRED_CALLBACK_URL")),
		facebook.New(os.Getenv("FACEBOOK_CRED_CRED_ID"), os.Getenv("FACEBOOK_CRED_SECRET"), os.Getenv("FACEBOOK_CRED_CALLBACK_URL")),
	)
}

type LoginOAuthUseCaseDeps struct {
	uberdig.In
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
	OauthConnectionRepository repository.OauthConnectionRepository
	TokenStrategy             token.TokenStrategy
}

type LoginOAuthUseCase struct {
	deps LoginOAuthUseCaseDeps
}

func NewLoginOAuthUseCase(deps LoginOAuthUseCaseDeps) *LoginOAuthUseCase {
	return &LoginOAuthUseCase{deps: deps}
}

// Execute handles the OAuth registration process, create a new user if it doesn't exist and generate tokens
func (u LoginOAuthUseCase) Execute(ctx context.Context, oauthUserId string) (*OAuthResult, error) {
	// Vérifier si l'utilisateur existe déjà
	oauthConnection, err := u.deps.OauthConnectionRepository.GetOauthConnectionByUserId(ctx, oauthUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth connection: %w", err)
	}

	user, err := u.deps.GetUserByIdUseCase.Execute(ctx, oauthConnection.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u.generateTokens(user, false)
}

type RegisterOAuthUseCaseDeps struct {
	uberdig.In
	GetUserByEmailUseCase *get_by_email.GetUserByEmailUseCase
	RegisterUserUseCase   *register.RegisterUserUseCase
}

type RegisterOAuthUseCase struct {
	deps RegisterOAuthUseCaseDeps
}

func NewRegisterOAuthUseCase(deps RegisterOAuthUseCaseDeps) *RegisterOAuthUseCase {
	return &RegisterOAuthUseCase{deps: deps}
}

func (u RegisterOAuthUseCase) Execute(ctx context.Context, provider string, oauthUserId string, oauthEmail string, inviteToken string) error {
	err := u.deps.RegisterUserUseCase.Execute(ctx, inviteToken, register.WithOauth(provider, oauthUserId, oauthEmail))
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

// Génération des tokens d'accès et de rafraîchissement
func (u LoginOAuthUseCase) generateTokens(user *entity.User, rememberMe bool) (*OAuthResult, error) {
	accessToken, err := u.deps.TokenStrategy.GenerateAccessToken(map[string]any{
		"email": user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := u.deps.TokenStrategy.GenerateRefreshToken(map[string]any{
		"email": user.Email,
	}, rememberMe)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &OAuthResult{
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenLifespan:  u.deps.TokenStrategy.GetAccessTokenLifespan(),
		RefreshTokenLifespan: u.deps.TokenStrategy.GetRefreshTokenLifespan(rememberMe),
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
