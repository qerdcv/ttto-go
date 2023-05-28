package xctx

import (
	"context"

	"github.com/qerdcv/ttto/internal/domain"
)

type ctxUserKey string

var userKey = ctxUserKey("user")

func ContextWithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) *domain.User {
	u, _ := ctx.Value(userKey).(*domain.User)
	return u
}
