package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/domain"
)

var (
	ErrInvalidJWTToken = errors.New("invalid jwt token")
)

type Claims struct {
	jwt.RegisteredClaims

	ID       int32  `json:"id"`
	Username string `json:"username"`
}

type JWTTokenizer struct {
	cfg conf.Auth
}

func NewJWTTokenizer(cfg conf.Auth) *JWTTokenizer {
	return &JWTTokenizer{
		cfg: cfg,
	}
}

func (t *JWTTokenizer) Encode(user *domain.User) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			ID:       user.ID,
			Username: user.Username,
		},
	).SignedString([]byte(t.cfg.Secret))
}

func (t *JWTTokenizer) Decode(accessToken string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signin method: %s", token.Header["alg"])
		}
		return []byte(t.cfg.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse with claims: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &domain.User{
			ID:       claims.ID,
			Username: claims.Username,
		}, nil
	}

	return nil, ErrInvalidJWTToken
}
