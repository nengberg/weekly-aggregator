package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Auth provides convenience functions for implementing the OAuth2 flow.
type Auth struct {
	Config *oauth2.Config
}

// NewAuth creates a new authenticator with configuration for OAuth2
func NewAuth() Auth {
	config := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"user-read-private", "playlist-modify-public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
	return Auth{
		Config: config,
	}
}

// AuthCodeURL ..
func (a *Auth) AuthCodeURL(state string) string {
	return a.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetToken bla
func (a *Auth) GetToken(r *http.Request) (*oauth2.Token, error) {
	token, err := a.Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Fatalf("Code exchange wrong: %s", err.Error())
	}
	return token, nil
}
