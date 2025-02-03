package token

import "time"

type TokenStrategy interface {
	GenerateAccessToken(claims map[string]any) (string, error)
	GenerateRefreshToken(claims map[string]any, longLived bool) (string, error)
	ValidateAccessToken(token string) (map[string]any, error)
	ValidateRefreshToken(token string) (map[string]any, error)
	GetAccessTokenLifespan() time.Duration
	GetRefreshTokenLifespan(longLived bool) time.Duration
}
