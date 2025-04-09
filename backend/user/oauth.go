package user

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/trysourcetool/sourcetool/backend/config"
)

var oauthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

const (
	googleOAuthCallbackPath = "/auth/google/authenticate"
)

type googleOAuthClient struct{}

func newGoogleOAuthClient() *googleOAuthClient {
	return &googleOAuthClient{}
}

type googleUserInfo struct {
	id         string
	email      string
	givenName  string
	familyName string
}

type googleToken struct {
	accessToken  string
	tokenType    string
	refreshToken string
	expiry       time.Time
}

func (c *googleOAuthClient) getGoogleAuthCodeURL(ctx context.Context, state string) (string, error) {
	redirectURL := config.Config.AuthBaseURL() + googleOAuthCallbackPath

	conf := &oauth2.Config{
		ClientID:     config.Config.Google.OAuth.ClientID,
		ClientSecret: config.Config.Google.OAuth.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       oauthScopes,
		Endpoint:     google.Endpoint,
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.ApprovalForce,
		oauth2.AccessTypeOffline,
	}

	return conf.AuthCodeURL(state, opts...), nil
}

func (c *googleOAuthClient) getGoogleToken(ctx context.Context, code string) (*googleToken, error) {
	redirectURL := config.Config.AuthBaseURL() + googleOAuthCallbackPath

	conf := &oauth2.Config{
		ClientID:     config.Config.Google.OAuth.ClientID,
		ClientSecret: config.Config.Google.OAuth.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       oauthScopes,
		Endpoint:     google.Endpoint,
	}

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return &googleToken{
		accessToken:  tok.AccessToken,
		tokenType:    tok.TokenType,
		refreshToken: tok.RefreshToken,
		expiry:       tok.Expiry,
	}, nil
}

func (c *googleOAuthClient) getGoogleUserInfo(ctx context.Context, tok *googleToken) (*googleUserInfo, error) {
	redirectURL := config.Config.AuthBaseURL() + googleOAuthCallbackPath

	conf := &oauth2.Config{
		ClientID:     config.Config.Google.OAuth.ClientID,
		ClientSecret: config.Config.Google.OAuth.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       oauthScopes,
		Endpoint:     google.Endpoint,
	}

	source := conf.TokenSource(ctx, &oauth2.Token{
		AccessToken:  tok.accessToken,
		TokenType:    tok.tokenType,
		RefreshToken: tok.refreshToken,
		Expiry:       tok.expiry,
	})

	service, err := googleOAuth2.NewService(ctx, option.WithTokenSource(source))
	if err != nil {
		return nil, err
	}

	info, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &googleUserInfo{
		id:         info.Id,
		email:      info.Email,
		givenName:  info.GivenName,
		familyName: info.FamilyName,
	}, nil
}
