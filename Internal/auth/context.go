package auth

import (
	"context"

	"github.com/google/uuid"
)

type Tenant struct {
	ID   uuid.UUID
	Plan string
}

type ctxKey struct{}

func WithTenant(ctx context.Context, t Tenant) context.Context {
	return context.WithValue(ctx, ctxKey{}, t)
}

// TenantFrom returns the tenant set by the auth middleware.
func TenantFrom(ctx context.Context) Tenant {
	t, _ := ctx.Value(ctxKey{}).(Tenant)
	return t
}
