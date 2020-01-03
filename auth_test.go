package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	clientID     = "123"
	clientSecret = "456"
	redirectURL  = "http://localhost:8080/callback"
)

func initializeClient() {
	os.Setenv("SPOTIFY_CLIENT_ID", clientID)
	os.Setenv("SPOTIFY_CLIENT_SECRET", clientSecret)
	os.Setenv("SPOTIFY_REDIRECT_URL", redirectURL)
}

func TestNewAuth(t *testing.T) {
	initializeClient()

	auth := NewAuth()

	assert.Equal(t, clientID, auth.Config.ClientID)
	assert.Equal(t, clientSecret, auth.Config.ClientSecret)
	assert.Equal(t, redirectURL, auth.Config.RedirectURL)
	assert.ElementsMatch(t, []string{"user-read-private", "playlist-modify-public"}, auth.Config.Scopes)
	assert.Equal(t, "https://accounts.spotify.com/authorize", auth.Config.Endpoint.AuthURL)
	assert.Equal(t, "https://accounts.spotify.com/api/token", auth.Config.Endpoint.TokenURL)
}

func TestAuthCodeURL(t *testing.T) {
	initializeClient()
	auth := NewAuth()

	assert.Equal(t,
		"https://accounts.spotify.com/authorize?"+
			"access_type=offline&client_id=123"+
			"&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback"+
			"&response_type=code&scope=user-read-private+playlist-modify-public"+
			"&state=123",
		auth.AuthCodeURL("123"),
	)
}
