package internal

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/encrypt"
)

func LoadFixtures(ctx context.Context, db database.DB) error {
	if !config.Config.IsCloudEdition {
		return nil
	}

	email := "john.doe@acme.com"
	exists, err := db.User().IsEmailExists(ctx, email)
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

	if err := db.WithTx(ctx, func(tx database.Tx) error {
		u := &core.User{
			ID:               uuid.Must(uuid.NewV4()),
			FirstName:        "John",
			LastName:         "Doe",
			Email:            email,
			RefreshTokenHash: hashedRefreshToken,
		}
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		o := &core.Organization{
			ID:        uuid.Must(uuid.NewV4()),
			Subdomain: internal.StringPtr("acme"),
		}
		if err := tx.Organization().Create(ctx, o); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, &core.UserOrganizationAccess{
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

		if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
			return err
		}

		key, err := core.GenerateAPIKey(devEnv.Slug)
		if err != nil {
			return err
		}

		keyHash, keyCiphertext, keyNonce, err := hashAndEncryptAPIKey(key)
		if err != nil {
			return err
		}
		apiKey := &core.APIKey{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			EnvironmentID:  devEnv.ID,
			UserID:         u.ID,
			Name:           "",
			KeyHash:        keyHash,
			KeyCiphertext:  keyCiphertext,
			KeyNonce:       keyNonce,
		}

		if err := tx.APIKey().Create(ctx, apiKey); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func hashAndEncryptAPIKey(plainAPIKey string) (keyHash string, keyCiphertext, keyNonce []byte, err error) {
	keyHash = core.HashAPIKey(plainAPIKey)

	encryptor, err := encrypt.NewEncryptor()
	if err != nil {
		return "", nil, nil, err
	}

	keyCiphertext, keyNonce, err = encryptor.Encrypt([]byte(plainAPIKey))
	if err != nil {
		return "", nil, nil, err
	}

	return keyHash, keyCiphertext, keyNonce, nil
}
