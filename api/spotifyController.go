package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SpotifyController struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

type spotifyTokenResponse struct {
	accessToken    string
	tokenType      string
	scope          string
	expirationTime int
	refreshToken   string
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

	if response.StatusCode != 200 {
		w.Write([]byte("Response error."))
	}
	if err != nil {
		log.Println("SpotifyAuthRequest %s", err.Error())
	}
}

func (sc *SpotifyController) AuthorizationCallback(w http.ResponseWriter, r *http.Request) {

	authCode := r.URL.Query()["code"][0]

	form := url.Values{}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))

	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("redirect_uri", sc.redirectURI)
	form.Add("client_id", sc.clientID)
	form.Add("client_secret", sc.clientSecret)
	req.PostForm = form

	if err != nil {
		log.Println(err)
	}

	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Println(readErr)
	}

	tokenResult := spotifyTokenResponse{}

	err = json.Unmarshal(body, &tokenResult)

	if err != nil {
		log.Println(err)
	}

	getAllUserArtists(tokenResult.accessToken)
}

func getAllUserArtists(userToken string) {
	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}
}
