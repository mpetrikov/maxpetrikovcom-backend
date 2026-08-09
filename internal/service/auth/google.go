package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config *oauth2.Config
}

func NewGoogleOAuth(
	clientID string,
	clientSecret string,
	redirectURL string,
) *GoogleOAuth {
	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
		},
	}
}

func (g *GoogleOAuth) AuthorizationURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func GenerateOAuthState() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
