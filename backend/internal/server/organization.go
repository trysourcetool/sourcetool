package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

type organizationResponse struct {
	ID                string `json:"id"`
	Subdomain         string `json:"subdomain"`
	WebSocketEndpoint string `json:"webSocketEndpoint"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

func organizationFromModel(o *core.Organization) *organizationResponse {
	if o == nil {
		return nil
	}

	return &organizationResponse{
		ID:                o.ID.String(),
		Subdomain:         internal.StringValue(o.Subdomain),
		WebSocketEndpoint: config.Config.WebSocketOrgBaseURL(internal.StringValue(o.Subdomain)),
		CreatedAt:         strconv.FormatInt(o.CreatedAt.Unix(), 10),
		UpdatedAt:         strconv.FormatInt(o.UpdatedAt.Unix(), 10),
	}
}

type createOrganizationRequest struct {
	Subdomain string `json:"subdomain" validate:"required"`
}

type createOrganizationResponse struct {
	Organization *organizationResponse `json:"organization"`
}

func (s *Server) handleCreateOrganization(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req createOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	var subdomain *string
	if config.Config.IsCloudEdition {
		subdomain = internal.StringPtr(req.Subdomain)

		if core.IsReservedSubdomain(req.Subdomain) {
			return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
		}

		exists, err := s.db.Organization().IsSubdomainExists(ctx, req.Subdomain)
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

	ctxUser := internal.ContextUser(ctx)

	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         ctxUser.ID,
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

	_, hashedKey, ciphertext, nonce, err := devEnv.GenerateAPIKey()
	if err != nil {
		return errdefs.ErrInternal(err)
	}
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         ctxUser.ID,
		Name:           "",
		KeyHash:        hashedKey,
		KeyCiphertext:  ciphertext,
		KeyNonce:       nonce,
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Organization().Create(ctx, o); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
			return err
		}

		if err := tx.APIKey().Create(ctx, apiKey); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	o, _ = s.db.Organization().Get(ctx, database.OrganizationByID(o.ID))

	return s.renderJSON(w, http.StatusOK, createOrganizationResponse{
		Organization: organizationFromModel(o),
	})
}

func (s *Server) handleCheckOrganizationSubdomainAvailability(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	subdomain := r.URL.Query().Get("subdomain")
	if subdomain == "" {
		return errdefs.ErrInvalidArgument(errors.New("subdomain is required"))
	}

	exists, err := s.db.Organization().IsSubdomainExists(ctx, subdomain)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
	}

	if core.IsReservedSubdomain(subdomain) {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
	}

	return s.renderJSON(w, http.StatusOK, statusResponse{
		Code:    http.StatusOK,
		Message: "Subdomain is available",
	})
}

func (s *Server) validateSelfHostedOrganization(ctx context.Context) error {
	if !config.Config.IsCloudEdition {
		// In self-hosted mode, check if an organization already exists
		if _, err := s.db.Organization().Get(ctx); err == nil {
			return errdefs.ErrPermissionDenied(errors.New("only one organization is allowed in self-hosted edition"))
		}
	}
	return nil
}

func (s *Server) createInitialOrganizationForSelfHosted(ctx context.Context, tx database.Tx, u *core.User) error {
	if config.Config.IsCloudEdition {
		return nil
	}

	org := &core.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: nil, // Empty subdomain for non-cloud edition
	}
	if err := tx.Organization().Create(ctx, org); err != nil {
		return err
	}

	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: org.ID,
		Role:           core.UserOrganizationRoleAdmin,
	}
	if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
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
	if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
		return err
	}

	_, hashedKey, ciphertext, nonce, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		KeyHash:        hashedKey,
		KeyCiphertext:  ciphertext,
		KeyNonce:       nonce,
	}
	if err := tx.APIKey().Create(ctx, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *Server) resolveOrganization(ctx context.Context, u *core.User) (*core.Organization, *core.UserOrganizationAccess, error) {
	if u == nil {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("user cannot be nil"))
	}

	ctxSubdomain := internal.ContextSubdomain(ctx)
	isCloudWithSubdomain := config.Config.IsCloudEdition && ctxSubdomain != "" && ctxSubdomain != "auth"

	orgAccessQueries := []database.UserOrganizationAccessQuery{
		database.UserOrganizationAccessByUserID(u.ID),
		database.UserOrganizationAccessOrderBy("created_at DESC"),
	}

	if isCloudWithSubdomain {
		orgAccessQueries = append(orgAccessQueries, database.UserOrganizationAccessByOrganizationSubdomain(ctxSubdomain))
	}

	orgAccess, err := s.db.User().GetOrganizationAccess(ctx, orgAccessQueries...)
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.db.Organization().Get(ctx, database.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}
