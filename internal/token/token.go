package token

import (
	"context"
	"fmt"
	"time"

	"github.com/aga-absolut/Vault-System/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// Provider defines JWT token operations.
//
//go:generate mockgen -source=token.go -destination=mocks/token_mock.go -package=mocks
type Provider interface {
	BuildJWTString(user string) (string, error)
	ValidateToken(tokenString string) (string, error)
	UserNameFromContext(ctx context.Context) (string, bool)
}

// UserIDKey is used as a context key for storing usernames.
type UserIDKey struct{}

// Claims represents custom JWT claims.
type Claims struct {
	jwt.RegisteredClaims
	Username string
}

// JWTProvider provides JWT token generation and validation.
type JWTProvider struct {
	secretKey []byte
	tokenTTL  time.Duration
}

// NewJWTProvider creates a new JWT provider instance.
func NewJWTProvider(cfg *config.Config) Provider {
	return &JWTProvider{
		secretKey: []byte(cfg.JWTSecret),
		tokenTTL:  cfg.TokenTTL,
	}
}

// BuildJWTString generates a signed JWT token for a user.
func (p *JWTProvider) BuildJWTString(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.tokenTTL)),
		},
		Username: user,
	})

	tokenString, err := token.SignedString([]byte(p.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the username.
func (p *JWTProvider) ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return p.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", err
	}
	return claims.Username, nil
}

// UserNameFromContext extracts the username from request context.
func (p *JWTProvider) UserNameFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(string)
	return userID, ok
}