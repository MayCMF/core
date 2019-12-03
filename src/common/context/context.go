package context

import (
	"context"
)

// Define keys in the global context
type (
	transCtx     struct{}
	transLockCtx struct{}
	userUUIDCtx  struct{}
	traceIDCtx   struct{}
)

// NewTrans - Create the context of the transaction
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// FromTrans - Get transactions from context
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// NewTransLock - Create a context for a transaction lock
func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, struct{}{})
}

// FromTransLock - Get transaction locks from context
func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil
}

// NewUserUUID - Create a context for the user ID
func NewUserUUID(ctx context.Context, userUUID string) context.Context {
	return context.WithValue(ctx, userUUIDCtx{}, userUUID)
}

// FromUserUUID - Get the user ID from the context
func FromUserUUID(ctx context.Context) (string, bool) {
	v := ctx.Value(userUUIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}

// NewTraceID - Create a context for tracking IDs
func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

// FromTraceID - Get the tracking ID from the context
func FromTraceID(ctx context.Context) (string, bool) {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}
