package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type reqIDContextKey struct{}

// GetReqID function    retrieves a request ID from given context.
func GetReqID(ctx context.Context) uuid.UUID {
	entry := ctx.Value(reqIDContextKey{})

	if reqID, ok := entry.(uuid.UUID); ok {
		return reqID
	}

	return uuid.Nil
}

// NewReqID function    creates a middleware for enrich context with unique request ID.
func NewReqID() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			reqID := uuid.New()
			next(context.WithValue(ctx, reqIDContextKey{}, reqID), b, update)
		}
	}
}
