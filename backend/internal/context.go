package internal

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type ctxKey string

const (
	ContextUserKey         ctxKey = "user"
	ContextOrganizationKey ctxKey = "organization"
	ContextSubdomainKey    ctxKey = "subdomain"
)

func ContextUser(ctx context.Context) *core.User {
	v, ok := ctx.Value(ContextUserKey).(*core.User)
	if !ok {
		return nil
	}
	return v
}

func ContextOrganization(ctx context.Context) *core.Organization {
	v, ok := ctx.Value(ContextOrganizationKey).(*core.Organization)
	if !ok {
		return nil
	}
	return v
}

func ContextSubdomain(ctx context.Context) string {
	v, ok := ctx.Value(ContextSubdomainKey).(string)
	if !ok {
		return ""
	}
	return v
}

func withUser(ctx context.Context, user *core.User) context.Context {
	return context.WithValue(ctx, ContextUserKey, user)
}

func withOrganization(ctx context.Context, org *core.Organization) context.Context {
	return context.WithValue(ctx, ContextOrganizationKey, org)
}

func withSubdomain(ctx context.Context, subdomain string) context.Context {
	return context.WithValue(ctx, ContextSubdomainKey, subdomain)
}

func NewBackgroundContext(ctx context.Context) context.Context {
	bgCtx := context.Background()

	if user := ContextUser(ctx); user != nil {
		bgCtx = withUser(bgCtx, user)
	}
	if org := ContextOrganization(ctx); org != nil {
		bgCtx = withOrganization(bgCtx, org)
	}
	if subdomain := ContextSubdomain(ctx); subdomain != "" {
		bgCtx = withSubdomain(bgCtx, subdomain)
	}

	return bgCtx
}
