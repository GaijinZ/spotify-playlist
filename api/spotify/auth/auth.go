package auth

import (
	"context"
	"fmt"
	"net/http"

	"spf-playlist/pkg/config"
	"spf-playlist/utils"

	"golang.org/x/oauth2"
)

type SpotifyAuth struct {
	env config.GlobalEnv
	ctx context.Context
}

func NewSpotifyAuth(env config.GlobalEnv, ctx context.Context) *SpotifyAuth {
	return &SpotifyAuth{
		env: env,
		ctx: ctx,
	}
}

// SpotifyAuth authorize user in Spotify API.
func (s *SpotifyAuth) SpotifyAuth(w http.ResponseWriter, r *http.Request) {
	log := utils.GetLogger(s.ctx)

	conf := &oauth2.Config{
		ClientID:     s.env.ClientID,
		ClientSecret: s.env.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  s.env.AuthURL,
			TokenURL: s.env.TokenURL,
		},
		RedirectURL: s.env.RedirectURI,
		Scopes:      []string{s.env.Scope},
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	log.Infof("Visit the URL for the auth dialog: %v\n", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// ExchangeToken exchanges the Spotify authorization code for an access token.
func (s *SpotifyAuth) ExchangeToken(code string) (*oauth2.Token, error) {
	log := utils.GetLogger(s.ctx)

	conf := &oauth2.Config{
		ClientID:     s.env.ClientID,
		ClientSecret: s.env.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: s.env.TokenURL,
		},
		RedirectURL: s.env.RedirectURI,
	}

	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Errorf("Failed to exchange token: %v", err)
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, nil
}
