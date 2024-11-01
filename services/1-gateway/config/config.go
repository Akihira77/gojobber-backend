package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig map[string]*oauth2.Config

func GetGoogleOAuthConfig(name string) (*oauth2.Config, error) {
	cfg, ok := GoogleOAuthConfig[name]
	if !ok {
		return nil, fmt.Errorf("OAuth Config For [%s] Action Did Not Exists", name)
	}

	return cfg, nil
}

func NewGoogleAuthConfig(port string) {
	GoogleOAuthConfig = make(map[string]*oauth2.Config)

	GoogleOAuthConfig["signin"] = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost%s/api/v1/gateway/auths/signin/google-callback", port),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	GoogleOAuthConfig["signup"] = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost%s/api/v1/gateway/auths/signup/google-callback", port),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
