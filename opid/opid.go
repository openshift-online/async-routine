package opid

import (
	"context"

	"github.com/google/uuid"
)

// ContextKey is the type of keys used to store operation identifiers in contexts.
type ContextKey int

const (
	// opidKey is the key used to store operation identifiers in context
	opidKey ContextKey = iota
)

// FromContext returns the operation identifier from the given context, or an empty string if no
// operation identifier is attached to the context.
func FromContext(ctx context.Context) string {
	return fromContext(ctx, opidKey)
}

func fromContext(ctx context.Context, key ContextKey) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

// IntoContext creates a new context containing the given operation identifier.
func IntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, opidKey, id)
}

// WithOpId creates a new context on top of the new one that contains a new operation identifier.
func WithOpId(ctx context.Context) context.Context {
	if ctx.Value(opidKey) != nil {
		return ctx
	}
	opId := uuid.NewString()
	return context.WithValue(ctx, opidKey, opId)
}

func NewContext() context.Context {
	ctx := context.Background()
	return WithOpId(ctx)
}

func CopyOpId(src context.Context, dst context.Context) context.Context {
	opID := src.Value(opidKey)
	return context.WithValue(dst, opidKey, opID)
}
