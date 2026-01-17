package middleware

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func SetUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int32, bool) {
	userID, ok := ctx.Value(userIDKey).(int32)
	return userID, ok
}
