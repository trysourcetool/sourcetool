package internal

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
)

type ctxKey string

const (
	CurrentUserCtxKey         ctxKey = "currentUser"
	CurrentOrganizationCtxKey ctxKey = "currentOrganization"
	SubdomainCtxKey           ctxKey = "subdomain"
)

func CurrentUser(ctx context.Context) *user.User {
	v, ok := ctx.Value(CurrentUserCtxKey).(*user.User)
	if !ok {
		return nil
	}
	return v
}

func CurrentOrganization(ctx context.Context) *organization.Organization {
	v, ok := ctx.Value(CurrentOrganizationCtxKey).(*organization.Organization)
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

func withUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, CurrentUserCtxKey, user)
}

func withOrganization(ctx context.Context, org *organization.Organization) context.Context {
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
