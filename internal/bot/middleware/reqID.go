package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

const reqIDKey middlewareCtxKey = "reqID"

func GetReqID(ctx context.Context) uuid.UUID {
	entry := ctx.Value(reqIDKey)

	if reqID, ok := entry.(uuid.UUID); ok {
		return reqID
	}

	return uuid.Nil
}

func NewReqID() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			reqID := uuid.New()
			next(context.WithValue(ctx, reqIDKey, reqID), b, update)
		}
	}
}
