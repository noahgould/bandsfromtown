package api

import (
	"log"
	"net/http"
	"os"
	"time"
)

type SpotifyController struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewSpotifyController() *SpotifyController {

	return &SpotifyController{
		clientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		redirectURI:  os.Getenv("SPOTIFY_REDIRECT_URL"),
	}
}

func (sc *SpotifyController) AuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	if err != nil {
		log.Println("SpotifyAuthRequest %s", err.Error())
	}

	q := req.URL.Query()
	q.Add("client_id", sc.clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", sc.redirectURI)
	q.Add("scope", "user-library-read playlist-read-collaborative playlist-read-private")

	req.URL.RawQuery = q.Encode()
	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := spotifyClient.Do(req)
	log.Println(req.RequestURI)
	log.Println(req.URL)
	log.Println(response.Status)

	if err != nil {
		log.Println("SpotifyAuthRequest %s", err.Error())
	}
}
