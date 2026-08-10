package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config *oauth2.Config
}

type GoogleUser struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

const googleUserInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

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
		return "", fmt.Errorf("generate oauth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (g *GoogleOAuth) GetUser(
	ctx context.Context,
	code string,
) (GoogleUser, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("exchange google oauth code: %w", err)
	}

	client := g.config.Client(ctx, token)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		googleUserInfoURL,
		nil,
	)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("create google userinfo request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("get google userinfo: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return GoogleUser{}, fmt.Errorf(
			"google userinfo returned status %d",
			response.StatusCode,
		)
	}

	var googleUser GoogleUser

	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		return GoogleUser{}, fmt.Errorf("decode google userinfo: %w", err)
	}

	return googleUser, nil
}
