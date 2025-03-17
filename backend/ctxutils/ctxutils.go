package ctxutils

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/model"
)

type ctxKey string

const (
	CurrentUserCtxKey         ctxKey = "currentUser"
	CurrentOrganizationCtxKey ctxKey = "currentOrganization"
	HTTPHostCtxKey            ctxKey = "host"
)

func CurrentUser(ctx context.Context) *model.User {
	v, ok := ctx.Value(CurrentUserCtxKey).(*model.User)
	if !ok {
		return nil
	}
	return v
}

func CurrentOrganization(ctx context.Context) *model.Organization {
	v, ok := ctx.Value(CurrentOrganizationCtxKey).(*model.Organization)
	if !ok {
		return nil
	}
	return v
}

func HTTPHost(ctx context.Context) string {
	v, ok := ctx.Value(HTTPHostCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func withUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, CurrentUserCtxKey, user)
}

func withOrganization(ctx context.Context, org *model.Organization) context.Context {
	return context.WithValue(ctx, CurrentOrganizationCtxKey, org)
}

func withHTTPHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, HTTPHostCtxKey, host)
}

func NewBackgroundContext(ctx context.Context) context.Context {
	bgCtx := context.Background()

	if user := CurrentUser(ctx); user != nil {
		bgCtx = withUser(bgCtx, user)
	}
	if org := CurrentOrganization(ctx); org != nil {
		bgCtx = withOrganization(bgCtx, org)
	}
	if host := HTTPHost(ctx); host != "" {
		bgCtx = withHTTPHost(bgCtx, host)
	}

	return bgCtx
}
