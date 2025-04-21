package internal

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type ctxKey string

const (
	CurrentUserCtxKey         ctxKey = "currentUser"
	CurrentOrganizationCtxKey ctxKey = "currentOrganization"
	SubdomainCtxKey           ctxKey = "subdomain"
)

func CurrentUser(ctx context.Context) *core.User {
	v, ok := ctx.Value(CurrentUserCtxKey).(*core.User)
	if !ok {
		return nil
	}
	return v
}

func CurrentOrganization(ctx context.Context) *core.Organization {
	v, ok := ctx.Value(CurrentOrganizationCtxKey).(*core.Organization)
	if !ok {
		return nil
	}
	return v
}

func Subdomain(ctx context.Context) string {
	v, ok := ctx.Value(SubdomainCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func withUser(ctx context.Context, user *core.User) context.Context {
	return context.WithValue(ctx, CurrentUserCtxKey, user)
}

func withOrganization(ctx context.Context, org *core.Organization) context.Context {
	return context.WithValue(ctx, CurrentOrganizationCtxKey, org)
}

func withSubdomain(ctx context.Context, subdomain string) context.Context {
	return context.WithValue(ctx, SubdomainCtxKey, subdomain)
}

func NewBackgroundContext(ctx context.Context) context.Context {
	bgCtx := context.Background()

	if user := CurrentUser(ctx); user != nil {
		bgCtx = withUser(bgCtx, user)
	}
	if org := CurrentOrganization(ctx); org != nil {
		bgCtx = withOrganization(bgCtx, org)
	}
	if subdomain := Subdomain(ctx); subdomain != "" {
		bgCtx = withSubdomain(bgCtx, subdomain)
	}

	return bgCtx
}
