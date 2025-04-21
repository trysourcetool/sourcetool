package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) createOrganization(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	var subdomain *string
	if config.Config.IsCloudEdition {
		subdomain = internal.NilValue(req.Subdomain)

		if lo.Contains(core.ReservedSubdomains, req.Subdomain) {
			return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
		}

		exists, err := s.db.IsOrganizationSubdomainExists(ctx, req.Subdomain)
		if err != nil {
			return err
		}
		if exists {
			return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
		}
	}

	o := &core.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: subdomain,
	}

	currentUser := internal.CurrentUser(ctx)

	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         currentUser.ID,
		OrganizationID: o.ID,
		Role:           core.UserOrganizationRoleAdmin,
	}
	devEnv := &core.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		Name:           core.EnvironmentNameDevelopment,
		Slug:           core.EnvironmentSlugDevelopment,
		Color:          core.EnvironmentColorDevelopment,
	}
	envs := []*core.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Name:           core.EnvironmentNameProduction,
			Slug:           core.EnvironmentSlugProduction,
			Color:          core.EnvironmentColorProduction,
		},
		devEnv,
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return errdefs.ErrInternal(err)
	}
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         currentUser.ID,
		Name:           "",
		Key:            key,
	}

	if err := s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		if err := s.db.CreateOrganization(ctx, tx, o); err != nil {
			return err
		}

		if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
			return err
		}

		if err := s.db.BulkInsertEnvironments(ctx, tx, envs); err != nil {
			return err
		}

		if err := s.db.CreateAPIKey(ctx, tx, apiKey); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	o, _ = s.db.GetOrganization(ctx, postgres.OrganizationByID(o.ID))

	return s.renderJSON(w, http.StatusOK, responses.CreateOrganizationResponse{
		Organization: responses.OrganizationFromModel(o),
	})
}

func (s *Server) checkOrganizationSubdomainAvailability(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	subdomain := r.URL.Query().Get("subdomain")
	if subdomain == "" {
		return errdefs.ErrInvalidArgument(errors.New("subdomain is required"))
	}

	exists, err := s.db.IsOrganizationSubdomainExists(ctx, subdomain)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
	}

	if lo.Contains(core.ReservedSubdomains, subdomain) {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
	}

	return s.renderJSON(w, http.StatusOK, responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Subdomain is available",
	})
}

// validateSelfHostedOrganization checks if creating a new organization is allowed in self-hosted mode.
func (s *Server) validateSelfHostedOrganization(ctx context.Context) error {
	if !config.Config.IsCloudEdition {
		// In self-hosted mode, check if an organization already exists
		if _, err := s.db.GetOrganization(ctx); err == nil {
			return errdefs.ErrPermissionDenied(errors.New("only one organization is allowed in self-hosted edition"))
		}
	}
	return nil
}

func (s *Server) createInitialOrganizationForSelfHosted(ctx context.Context, tx *sqlx.Tx, u *core.User) error {
	if config.Config.IsCloudEdition {
		return nil
	}

	org := &core.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: nil, // Empty subdomain for non-cloud edition
	}
	if err := s.db.CreateOrganization(ctx, tx, org); err != nil {
		return err
	}

	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: org.ID,
		Role:           core.UserOrganizationRoleAdmin,
	}
	if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
		return err
	}

	devEnv := &core.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		Name:           core.EnvironmentNameDevelopment,
		Slug:           core.EnvironmentSlugDevelopment,
		Color:          core.EnvironmentColorDevelopment,
	}
	envs := []*core.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: org.ID,
			Name:           core.EnvironmentNameProduction,
			Slug:           core.EnvironmentSlugProduction,
			Color:          core.EnvironmentColorProduction,
		},
		devEnv,
	}
	if err := s.db.BulkInsertEnvironments(ctx, tx, envs); err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		Key:            key,
	}
	if err := s.db.CreateAPIKey(ctx, tx, apiKey); err != nil {
		return err
	}

	return nil
}

// resolveOrganizationBySubdomain gets an organization by subdomain and verifies the user has access.
// Deprecated: Use getOrganizationBySubdomain instead.
func (s *Server) resolveOrganizationBySubdomain(ctx context.Context, u *core.User, subdomain string) (*core.Organization, *core.UserOrganizationAccess, error) {
	if subdomain == "" {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("subdomain cannot be empty"))
	}

	return s.getOrganizationBySubdomain(ctx, u, subdomain)
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *Server) getUserOrganizationInfo(ctx context.Context) (*core.Organization, *core.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, internal.CurrentUser(ctx))
}

// getOrganizationBySubdomain retrieves an organization by subdomain and verifies user access.
func (s *Server) getOrganizationBySubdomain(ctx context.Context, u *core.User, subdomain string) (*core.Organization, *core.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.db.GetUserOrganizationAccess(ctx,
		postgres.UserOrganizationAccessByOrganizationID(org.ID),
		postgres.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *Server) getOrganizationInfo(ctx context.Context, u *core.User) (*core.Organization, *core.UserOrganizationAccess, error) {
	if u == nil {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("user cannot be nil"))
	}

	subdomain := internal.Subdomain(ctx)
	isCloudWithSubdomain := config.Config.IsCloudEdition && subdomain != "" && subdomain != "auth"

	// Different strategies for cloud vs. self-hosted or auth subdomain
	if isCloudWithSubdomain {
		return s.getOrganizationBySubdomain(ctx, u, subdomain)
	}

	return s.getDefaultOrganizationForUser(ctx, u)
}

// (typically the most recently created one).
func (s *Server) getDefaultOrganizationForUser(ctx context.Context, u *core.User) (*core.Organization, *core.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.db.GetUserOrganizationAccess(ctx,
		postgres.UserOrganizationAccessByUserID(u.ID),
		postgres.UserOrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}
