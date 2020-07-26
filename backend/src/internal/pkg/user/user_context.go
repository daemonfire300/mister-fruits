package user

import "context"

type UserContextKey string

const contextKey UserContextKey = "userctx"

func FromContext(ctx context.Context) string {
	username := ctx.Value(contextKey)
	// TODO should handle more cases, i.e., if the casting attempt below fails
	return username.(string)
}

func ContextWithUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, contextKey, username)
}
