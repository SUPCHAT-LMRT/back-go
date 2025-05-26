package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	getUserByEmailUseCase *get_by_email.GetUserByEmailUseCase
	tokenStrategy         token.TokenStrategy
}

var (
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrUserNotVerified           = errors.New("user not verified")
)

func NewLoginUserUseCase(
	getUserByEmailUseCase *get_by_email.GetUserByEmailUseCase,
	tokenStrategy token.TokenStrategy,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		getUserByEmailUseCase: getUserByEmailUseCase,
		tokenStrategy:         tokenStrategy,
	}
}

func (l *LoginUserUseCase) Execute(
	ctx context.Context,
	request LoginUserRequest,
) (*LoginUserResult, error) {
	user, err := l.getUserByEmailUseCase.Execute(
		ctx,
		request.Email,
		get_by_email.WithUserPassword(),
	)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return nil, ErrInvalidUsernameOrPassword
	}

	accessToken, err := l.tokenStrategy.GenerateAccessToken(map[string]any{
		"email": request.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	refreshToken, err := l.tokenStrategy.GenerateRefreshToken(map[string]any{
		"email": request.Email,
	}, request.RememberMe)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &LoginUserResult{
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenLifespan:  l.tokenStrategy.GetAccessTokenLifespan(),
		RefreshTokenLifespan: l.tokenStrategy.GetRefreshTokenLifespan(request.RememberMe),
		User:                 user,
	}, nil
}

type LoginUserRequest struct {
	Email      string
	Password   string
	RememberMe bool
}

type LoginUserResult struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenLifespan  time.Duration
	RefreshTokenLifespan time.Duration
	User                 *entity.User
}
