package context

import (
	"context"

	"github.com/msa-ali/picbucket/models"
)

type key int

const (
	userKey key = iota
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(context.Background(), userKey, user)
}

func User(ctx context.Context) *models.User {
	user, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
