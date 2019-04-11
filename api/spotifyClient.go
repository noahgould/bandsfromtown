package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type SpotifyClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func NewSpotifyClient() *SpotifyClient {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	clientID, clientEnvExist := os.LookupEnv("SPOTIFY_ID")
	clientSecret, clientSecretEnvExist := os.LookupEnv("SPOTIFY_SECRET")
	redirectURL, redirectURLEnvExist := os.LookupEnv("SPOTIFY_REDIRECT_URL")
	if !clientEnvExist {
		log.Fatal("spotify client id  not stored in environment variables.")
	}
	if !clientSecretEnvExist {
		log.Fatal("spotify client secret not stored in environment variables.")
	}
	if !redirectURLEnvExist {
		log.Fatal("spotify redirect url not stored in environment variables.")
	}

	return &SpotifyClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (sc *SpotifyClient) startSpotifySession(authCode string) spotifyTokenResponse {

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("redirect_uri", sc.redirectURL)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println(err)
	}

	headerString := base64.StdEncoding.EncodeToString([]byte(sc.clientID + ":" + sc.clientSecret))
	headerString = strings.Join([]string{"Basic", headerString}, " ")
	req.Header.Add("Authorization", headerString)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	spotifyHttpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := spotifyHttpClient.Do(req)

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

	return tokenResult
}

func makeAlbumRequest(userToken string, offset int) spotifyAlbumPage {

	spotifyHttpClient := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/albums", nil)

	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyHttpClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyAlbumPage{}

	if response.StatusCode == 200 {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		jsonErr := json.Unmarshal(body, &firstPage)

		if jsonErr != nil {
			log.Println(jsonErr)
		}
	}

	return firstPage
}

func makePlaylistRequest(userToken string, offset int) spotifyPlaylistPage {
	spotifyHttpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)

	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyHttpClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyPlaylistPage{}

	if response.StatusCode == 200 {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		jsonErr := json.Unmarshal(body, &firstPage)

		if jsonErr != nil {
			log.Println(jsonErr)
		}
	}

	return firstPage
}

// can make this more efficient by limiting results. https://beta.developer.spotify.com/documentation/web-api/reference/playlists/get-playlists-tracks/
func makePlaylistTrackRequest(userToken string, offset int, playlist spotifySimplePlaylist) spotifyTrackPage {

	spotifyHttpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", playlist.Owner.ID, playlist.ID)

	req, err := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "100")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyHttpClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyTrackPage{}

	if response.StatusCode == 200 {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Println(readErr)
		}
		jsonErr := json.Unmarshal(body, &firstPage)

		if jsonErr != nil {
			log.Println(jsonErr)
		}
	}

	return firstPage
}
