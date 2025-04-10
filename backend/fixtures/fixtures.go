package fixtures

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
)

func Load(ctx context.Context, store infra.Store) error {
	if !config.Config.IsCloudEdition {
		return nil
	}

	email := "john.doe@acme.com"
	exists, err := store.User().IsEmailExists(ctx, email)
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

	return store.RunTransaction(func(tx infra.Transaction) error {
		u := &model.User{
			ID:               uuid.Must(uuid.NewV4()),
			FirstName:        "John",
			LastName:         "Doe",
			Email:            email,
			RefreshTokenHash: hashedRefreshToken,
		}
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		o := &model.Organization{
			ID:        uuid.Must(uuid.NewV4()),
			Subdomain: conv.NilValue("acme"),
		}
		if err := tx.Organization().Create(ctx, o); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, &model.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: o.ID,
			Role:           model.UserOrganizationRoleAdmin,
		}); err != nil {
			return err
		}

		devEnv := &model.Environment{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Name:           model.EnvironmentNameDevelopment,
			Slug:           model.EnvironmentSlugDevelopment,
			Color:          model.EnvironmentColorDevelopment,
		}
		envs := []*model.Environment{
			{
				ID:             uuid.Must(uuid.NewV4()),
				OrganizationID: o.ID,
				Name:           model.EnvironmentNameProduction,
				Slug:           model.EnvironmentSlugProduction,
				Color:          model.EnvironmentColorProduction,
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
		apiKey := &model.APIKey{
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
