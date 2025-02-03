package token

import "time"

type RefreshAccessTokenUseCase struct {
	tokenStrategy TokenStrategy
}

func NewRefreshAccessTokenUseCase(tokenStrategy TokenStrategy) *RefreshAccessTokenUseCase {
	return &RefreshAccessTokenUseCase{tokenStrategy: tokenStrategy}
}

func (u *RefreshAccessTokenUseCase) Execute(refreshToken string) (*RefreshAccessTokenResponse, error) {
	parsedRefreshToken, err := u.tokenStrategy.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.tokenStrategy.GenerateAccessToken(map[string]any{
		"email": parsedRefreshToken["email"],
	})
	if err != nil {
		return nil, err
	}

	return &RefreshAccessTokenResponse{
		AccessToken:         accessToken,
		AccessTokenLifespan: u.tokenStrategy.GetAccessTokenLifespan(),
	}, nil
}

type RefreshAccessTokenResponse struct {
	AccessToken         string
	AccessTokenLifespan time.Duration
}
