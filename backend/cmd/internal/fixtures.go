package internal

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
)

func LoadFixtures(ctx context.Context, db *postgres.DB) error {
	if !config.Config.IsCloudEdition {
		return nil
	}

	email := "john.doe@acme.com"
	exists, err := db.IsUserEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return err
	}

	if err := db.WithTx(ctx, func(tx *sqlx.Tx) error {
		u := &core.User{
			ID:               uuid.Must(uuid.NewV4()),
			FirstName:        "John",
			LastName:         "Doe",
			Email:            email,
			RefreshTokenHash: hashedRefreshToken,
		}
		if err := db.CreateUser(ctx, tx, u); err != nil {
			return err
		}

		o := &core.Organization{
			ID:        uuid.Must(uuid.NewV4()),
			Subdomain: internal.NilValue("acme"),
		}
		if err := db.CreateOrganization(ctx, tx, o); err != nil {
			return err
		}

		if err := db.CreateUserOrganizationAccess(ctx, tx, &core.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: o.ID,
			Role:           core.UserOrganizationRoleAdmin,
		}); err != nil {
			return err
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

		if err := db.BulkInsertEnvironments(ctx, tx, envs); err != nil {
			return err
		}

		key, err := devEnv.GenerateAPIKey()
		if err != nil {
			return err
		}
		apiKey := &core.APIKey{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			EnvironmentID:  devEnv.ID,
			UserID:         u.ID,
			Name:           "",
			Key:            key,
		}

		if err := db.CreateAPIKey(ctx, tx, apiKey); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
