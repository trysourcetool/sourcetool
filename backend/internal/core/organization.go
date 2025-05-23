package core

import (
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Subdomain *string   `db:"subdomain"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

var ReservedSubdomains = []string{
	// Brand Protection
	"sourcetool",
	"trysourcetool",

	// Authentication & Security
	"2fa",
	"auth",
	"auth0",
	"login",
	"mfa",
	"oauth",
	"password",
	"passwords",
	"register",
	"registration",
	"saml",
	"secure",
	"security",
	"signin",
	"signup",
	"sso",

	// Core Infrastructure
	"api",
	"api-docs",
	"apis",
	"app",
	"apps",
	"cache",
	"cdn",
	"config",
	"db",
	"database",
	"git",
	"graphql",
	"logs",
	"proxy",
	"repo",
	"repository",
	"server",
	"service",
	"services",
	"socket",
	"static",
	"storage",
	"svc",
	"sys",
	"system",
	"webhook",
	"webhooks",
	"ws",

	// Environment & Testing
	"alpha",
	"beta",
	"demo",
	"example",
	"examples",
	"playground",
	"preview",
	"prod",
	"production",
	"prd",
	"qa",
	"quality",
	"quality-assurance",
	"release",
	"sample",
	"samples",
	"sandbox",
	"staging",
	"stg",
	"test",
	"testing",
	"tests",

	// User & Organization Management
	"account",
	"accounts",
	"admin",
	"administrator",
	"groups",
	"members",
	"my",
	"org",
	"organization",
	"portal",
	"profile",
	"root",
	"settings",
	"super",
	"superuser",
	"team",
	"teams",
	"user",
	"users",
	"customer",
	"customers",
	"welcome",
	"workspace",
	"workspaces",
	"project",
	"projects",

	// Business & Enterprise
	"billing",
	"corporate",
	"enterprise",
	"finance",
	"hr",
	"invoice",
	"marketing",
	"payments",
	"premium",
	"pro",
	"sales",

	// Product Features
	"alerts",
	"analytics",
	"console",
	"dashboard",
	"feedback",
	"internal",
	"management",
	"metrics",
	"monitor",
	"notifications",
	"private",
	"public",
	"reports",
	"search",

	// Communication & Support
	"email",
	"ftp",
	"help",
	"imap",
	"mail",
	"pop",
	"smtp",
	"ssh",
	"support",

	// Documentation & Development
	"dev",
	"developer",
	"developers",
	"docs",
	"document",
	"documentation",
	"documentations",
	"documents",
	"swagger",

	// Content & Media
	"assets",
	"blog",
	"browser",
	"media",
	"news",
	"web",
	"www",

	// Legal & Company
	"about",
	"contact",
	"legal",
	"privacy",
	"terms",

	// Commerce
	"shop",
	"store",

	// Temporary & Misc
	"hello",
	"official",
	"temp",
	"temporary",
	"tmp",
}

func IsReservedSubdomain(subdomain string) bool {
	return slices.Contains(ReservedSubdomains, subdomain)
}
