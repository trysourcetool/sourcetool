package server

import (
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
