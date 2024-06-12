package service

import (
	"context"
	"encoding/json"

	"github.com/Daka-0424/my-go-server/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

type IOauthGoogle interface {
	OauthGoogleURL(ctx context.Context, endPoint string) string
	OauthGoogle(ctx context.Context, code, endPoint string) (*UserInfo, error)
}

type oauthGoogleService struct {
	cfg         *config.Config
	oauthConfig *oauth2.Config
}

func NewOauthGoogleService(
	cfg *config.Config,
) IOauthGoogle {
	service := &oauthGoogleService{
		cfg: cfg,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.Oauth.GoogleOauth.ClientID,
			ClientSecret: cfg.Oauth.GoogleOauth.ClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
	return service
}

func (service *oauthGoogleService) OauthGoogleURL(ctx context.Context, endPoint string) string {
	service.oauthConfig.RedirectURL = service.cfg.Oauth.GoogleOauth.RedirectURL + endPoint
	return service.oauthConfig.AuthCodeURL(service.cfg.Oauth.GoogleOauth.OauthStateString)
}

func (service *oauthGoogleService) OauthGoogle(ctx context.Context, code, endPoint string) (*UserInfo, error) {
	service.oauthConfig.RedirectURL = service.cfg.Oauth.GoogleOauth.RedirectURL + endPoint
	token, err := service.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	client := service.oauthConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
