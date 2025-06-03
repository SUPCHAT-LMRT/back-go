package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supchat-lmrt/back-go/internal/utils"
)

const (
	JwtAccessTokenLifespan           = 15 * time.Minute
	JwtRefreshTokenLifespan          = 24 * time.Hour
	JwtRefreshTokenLongLivedLifespan = 30 * 24 * time.Hour
)

var (
	ErrUnableToSignAccessToken      = fmt.Errorf("unable to sign access token")
	ErrUnableToSignRefreshToken     = fmt.Errorf("unable to sign refresh token")
	ErrUnableToValidateAccessToken  = fmt.Errorf("unable to validate access token")
	ErrAccessTokenExpired           = fmt.Errorf("access token expired")
	ErrUnableToValidateRefreshToken = fmt.Errorf("unable to validate refresh token")
	ErrRefreshTokenExpired          = fmt.Errorf("refresh token expired")
	ErrInvalidTokenType             = fmt.Errorf("invalid token type")
)

type JwtTokenStrategy struct {
	secret string
}

func NewJwtTokenStrategy(secret string) func() TokenStrategy {
	return func() TokenStrategy {
		return &JwtTokenStrategy{secret}
	}
}

func (j *JwtTokenStrategy) GenerateAccessToken(claims map[string]any) (string, error) {
	if claims == nil {
		claims = make(map[string]any)
	}
	claims["exp"] = time.Now().Add(JwtAccessTokenLifespan).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	signedAccessToken, err := accessToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnableToSignAccessToken, err)
	}

	return signedAccessToken, nil
}

func (j *JwtTokenStrategy) GenerateRefreshToken(
	claims map[string]any,
	longLived bool,
) (string, error) {
	if claims == nil {
		claims = make(map[string]any)
	}
	claims["exp"] = time.Now().
		Add(utils.IfThenElse(longLived, JwtRefreshTokenLongLivedLifespan, JwtRefreshTokenLifespan)).
		Unix()
	claims["type"] = "refresh"
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	signedRefreshToken, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnableToSignRefreshToken, err)
	}

	return signedRefreshToken, nil
}

func (j *JwtTokenStrategy) ValidateAccessToken(token string) (map[string]any, error) {
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnableToValidateAccessToken, err)
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok || !accessToken.Valid {
		return nil, fmt.Errorf("%w: %w", ErrAccessTokenExpired, err)
	}

	// If the token is not an access token, return an error
	if claims["type"] == "refresh" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func (j *JwtTokenStrategy) ValidateRefreshToken(token string) (map[string]any, error) {
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnableToValidateRefreshToken, err)
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		return nil, fmt.Errorf("%w: %w", ErrRefreshTokenExpired, err)
	}

	if claims["type"] != "refresh" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func (j *JwtTokenStrategy) GetAccessTokenLifespan() time.Duration {
	return JwtAccessTokenLifespan
}

func (j *JwtTokenStrategy) GetRefreshTokenLifespan(longLived bool) time.Duration {
	return utils.IfThenElse(longLived, JwtRefreshTokenLongLivedLifespan, JwtRefreshTokenLifespan)
}
