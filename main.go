package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron"
)

var (
	auth              = NewAuth()
	listToAddTracksTo = os.Getenv("SPOTIFY_AGGREGATION_LIST_ID")
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", authenticate)
	mux.HandleFunc("/callback", callback)

	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: mux,
	}

	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}

}

func authenticate(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthCodeURL("123")
	http.Redirect(w, r, url, 302)
}

func callback(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetToken(r)
	if err != nil {
		log.Fatalf("Couldn't get token: %s", err.Error())
	}
	client := auth.NewClient(token)

	c := cron.New()
	c.AddFunc("@every 0h10m0s", func() {
		playlists, _ := client.GetPlaylists()
		if err != nil {
			log.Fatalf("Error getting playlists: %s", err.Error())
		} else {
			for _, playlist := range playlists.Items {
				if playlist.Name == "Discover Weekly" {
					tracks, _ := client.GetTracks(playlist.GetID())
					client.AddTracksToPlaylist(listToAddTracksTo, tracks)
				}
			}
		}
	})
	c.Start()
}
