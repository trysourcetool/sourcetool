package fixtures

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/trysourcetool/sourcetool/backend/conv"
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

	_, hashedSecret, err := generateSecret()
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
			Secret:               hashedSecret,
			EmailAuthenticatedAt: &now,
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

		key, err := devEnv.GenerateAPIKey(conv.SafeValue(o.Subdomain))
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

func generateSecret() (plainSecret, hashedSecret string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainSecret = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainSecret))
	hashedSecret = hex.EncodeToString(hash[:])

	return plainSecret, hashedSecret, nil
}
