package fixtures

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/pkg/conv"
)

func Load(ctx context.Context, repo db.Repository) error {
	if !config.Config.IsCloudEdition {
		return nil
	}

	email := "john.doe@acme.com"
	exists, err := repo.User().IsEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return err
	}

	return repo.RunTransaction(func(tx db.Transaction) error {
		u := &user.User{
			ID:               uuid.Must(uuid.NewV4()),
			FirstName:        "John",
			LastName:         "Doe",
			Email:            email,
			RefreshTokenHash: hashedRefreshToken,
		}
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		o := &organization.Organization{
			ID:        uuid.Must(uuid.NewV4()),
			Subdomain: conv.NilValue("acme"),
		}
		if err := tx.Organization().Create(ctx, o); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, &user.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: o.ID,
			Role:           user.UserOrganizationRoleAdmin,
		}); err != nil {
			return err
		}

		devEnv := &environment.Environment{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Name:           environment.EnvironmentNameDevelopment,
			Slug:           environment.EnvironmentSlugDevelopment,
			Color:          environment.EnvironmentColorDevelopment,
		}
		envs := []*environment.Environment{
			{
				ID:             uuid.Must(uuid.NewV4()),
				OrganizationID: o.ID,
				Name:           environment.EnvironmentNameProduction,
				Slug:           environment.EnvironmentSlugProduction,
				Color:          environment.EnvironmentColorProduction,
			},
			devEnv,
		}

		if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
			return err
		}

		key, err := devEnv.GenerateAPIKey()
		if err != nil {
			return err
		}
		apiKey := &apikey.APIKey{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			EnvironmentID:  devEnv.ID,
			UserID:         u.ID,
			Name:           "",
			Key:            key,
		}

		if err := tx.APIKey().Create(ctx, apiKey); err != nil {
			return err
		}

		return nil
	})
}

func generateRefreshToken() (plainRefreshToken, hashedRefreshToken string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainRefreshToken = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainRefreshToken))
	hashedRefreshToken = hex.EncodeToString(hash[:])

	return plainRefreshToken, hashedRefreshToken, nil
}
