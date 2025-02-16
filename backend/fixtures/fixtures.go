package fixtures

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

func Load(ctx context.Context, store infra.Store) error {
	email := "john.doe@acme.com"
	exists, err := store.User().IsEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	encodedPass, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		return err
	}

	secret, err := generateSecret()
	if err != nil {
		return err
	}

	now := time.Now()
	return store.RunTransaction(func(tx infra.Transaction) error {
		u := &model.User{
			ID:                   uuid.Must(uuid.NewV4()),
			FirstName:            "John",
			LastName:             "Doe",
			Email:                email,
			Password:             hex.EncodeToString(encodedPass[:]),
			Secret:               secret,
			EmailAuthenticatedAt: &now,
		}
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		o := &model.Organization{
			ID:        uuid.Must(uuid.NewV4()),
			Subdomain: "acme",
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

		key, err := model.GenerateAPIKey(o.Subdomain, devEnv.Slug)
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

func generateSecret() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
