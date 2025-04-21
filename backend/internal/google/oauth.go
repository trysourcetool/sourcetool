package google

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

var oauthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

const (
	googleOAuthCallbackPath = "/auth/google/callback"
)

type OAuthClient struct{}

func NewOAuthClient() *OAuthClient {
	return &OAuthClient{}
}

type UserInfo struct {
	ID         string
	Email      string
	GivenName  string
	FamilyName string
}

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func (c *OAuthClient) GetGoogleAuthCodeURL(ctx context.Context, state string) (string, error) {
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

func (c *OAuthClient) GetGoogleToken(ctx context.Context, code string) (*Token, error) {
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

	return &Token{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}, nil
}

func (c *OAuthClient) GetGoogleUserInfo(ctx context.Context, tok *Token) (*UserInfo, error) {
	redirectURL := config.Config.AuthBaseURL() + googleOAuthCallbackPath

	conf := &oauth2.Config{
		ClientID:     config.Config.Google.OAuth.ClientID,
		ClientSecret: config.Config.Google.OAuth.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       oauthScopes,
		Endpoint:     google.Endpoint,
	}

	source := conf.TokenSource(ctx, &oauth2.Token{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	})

	service, err := googleOAuth2.NewService(ctx, option.WithTokenSource(source))
	if err != nil {
		return nil, err
	}

	info, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:         info.Id,
		Email:      info.Email,
		GivenName:  info.GivenName,
		FamilyName: info.FamilyName,
	}, nil
}
